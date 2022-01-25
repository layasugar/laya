package grpcx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	pbrpc "github.com/layasugar/laya/grpcx/pbrpc"
	"io"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
)

// error log info definition
var ERR_IGNORE_ERR = errors.New("[marshal-001]Ingore error")
var ERR_NO_SNAPPY = errors.New("[marshal-002]Snappy compress not support yet.")
var ERR_META = errors.New("[marshal-003]Get nil value from Meta struct after marshal")

/*
 Data package for baidu RPC.
 all request and response data package should apply this.

-----------------------------------
| Head | Meta | Data | Attachment |
-----------------------------------

1. <Head> with fixed 12 byte length as follow format
----------------------------------------------
| PRPC | MessageSize(int32) | MetaSize(int32) |
----------------------------------------------
MessageSize = totalSize - 12(Fixed Head Size)
MetaSize = Meta object size

2. <Meta> body proto description as follow
message RpcMeta {
    optional RpcRequestMeta request = 1;
    optional RpcResponseMeta response = 2;
    optional int32 compress_type = 3; // 0:nocompress 1:Snappy 2:gzip
    optional int64 correlation_id = 4;
    optional int32 attachment_size = 5;
    optional ChunkInfo chuck_info = 6;
    optional bytes authentication_data = 7;
};

message Request {
    required string service_name = 1;
    required string method_name = 2;
    optional int64 log_id = 3;
};

message Response {
    optional int32 error_code = 1;
    optional string error_text = 2;
};

messsage ChunkInfo {
        required int64 stream_id = 1;
        required int64 chunk_id = 2;
};

3. <Data> customize transport data message.

4. <Attachment> attachment body data message
*/
type Package struct {
	Header     Header
	Meta       pbrpc.RpcMeta
	Data       []byte
	Attachment []byte
}

func NewPackage() *Package {
	pkg := &Package{}
	pkg.Meta.Request = &pbrpc.Request{}
	pkg.Meta.Response = &pbrpc.Response{}
	return pkg
}

func NewRequestPackage() *Package {
	pkg := &Package{}
	pkg.Meta.Request = &pbrpc.Request{}
	return pkg
}

func NewResponsePackage() *Package {
	pkg := &Package{}
	pkg.Meta.Response = &pbrpc.Response{}
	return pkg
}

func (r *Package) SetMagicCode(magicCode []byte) error {
	return r.Header.SetMagicCode(magicCode)
}

func (r *Package) GetMagicCode() string {
	return string(r.Header.MagicCode)
}

// TODO: 看看协议，此处叫data还是叫payload合适？
func (r *Package) SetData(Data []byte) *Package {
	r.Data = Data
	return r
}

func (r *Package) SetAttachment(Attachment []byte) *Package {
	r.Attachment = Attachment
	return r
}

func (r *Package) SetServiceName(serviceName string) *Package {
	r.Meta.Request.ServiceName = *proto.String(serviceName)
	return r
}

func (r *Package) SetMethodName(methodName string) *Package {
	r.Meta.Request.MethodName = *proto.String(methodName)
	return r
}

func (r *Package) SetTraceId(traceId string) *Package {
	r.Meta.Request.TraceId = *proto.String(traceId)
	return r
}

func (r *Package) GetTraceId() string {
	return r.Meta.Request.GetTraceId()
}

func (r *Package) SetCorrelationId(correlationId int64) *Package {
	r.Meta.CorrelationId = proto.Int64(correlationId)
	return r
}

func (r *Package) SetCompressType(compressType int32) *Package {
	r.Meta.CompressType = proto.Int32(compressType)
	return r
}

func (r *Package) SetAuthenticationData(authenticationData []byte) *Package {
	r.Meta.AuthenticationData = authenticationData
	return r
}

func (r *Package) SetErrorCode(errorCode int32) *Package {
	r.Meta.Response.ErrorCode = proto.Int32(errorCode)
	return r
}

func (r *Package) SetErrorText(errorText string) *Package {
	r.Meta.Response.ErrorMsg = proto.String(errorText)
	return r
}

func (r *Package) SetChunkInfo(streamId int64, chunkId int64) *Package {
	r.Meta.ChunkInfo = &pbrpc.ChunkInfo{
		StreamId: *proto.Int64(streamId),
		ChunkId:  *proto.Int64(chunkId),
	}
	return r
}

// WriteIO write package to io.Writer
func (r *Package) WriteIO(rw io.Writer) (int, error) {
	data, err := r.Bytes()
	if err != nil {
		return 0, err
	}

	_, err = rw.Write(data)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

// Bytes Convert RpcPackage to byte array
func (r *Package) Bytes() ([]byte, error) {
	var err error
	var totalSize int32

	if r.Data != nil {
		switch r.Meta.GetCompressType() {
		case COMPRESS_GZIP:
			r.Data, err = GZIP(r.Data)
			if err != nil {
				return nil, err
			}
		case COMPRESS_SNAPPY:
			dst := make([]byte, snappy.MaxEncodedLen(len(r.Data)))
			r.Data = snappy.Encode(dst, r.Data)
		}

		totalSize += int32(len(r.Data))
	}

	var attachmentSize int32
	if r.Attachment != nil {
		attachmentSize = int32(len(r.Attachment))
		totalSize += attachmentSize
	}
	r.Meta.AttachmentSize = proto.Int32(int32(attachmentSize))

	metaBytes, err := proto.Marshal(&r.Meta)
	if err != nil {
		return nil, err
	}

	if metaBytes == nil {
		return nil, ERR_META
	}

	totalSize += int32(len(metaBytes))

	r.Header.MetaSize = int32(len(metaBytes))
	r.Header.MessageSize = totalSize // set message body size

	buf := new(bytes.Buffer)

	headBytes := r.Header.Bytes()
	binary.Write(buf, binary.BigEndian, headBytes)
	binary.Write(buf, binary.BigEndian, metaBytes)

	if r.Data != nil {
		binary.Write(buf, binary.BigEndian, r.Data)
	}

	if r.Attachment != nil {
		binary.Write(buf, binary.BigEndian, r.Attachment)
	}

	return buf.Bytes(), nil
}

// ReadIO Read byte array and initialize RpcPackage
func (r *Package) ReadIO(rw io.Reader) error {
	if rw == nil {
		return errors.New("bytes is nil")
	}

	// read Head
	header := make([]byte, HeaderSize)
	_, err := io.ReadFull(rw, header)
	if err != nil {
		log.Println("Read head error", err)
		// only to close current connection
		return ERR_IGNORE_ERR
	}

	// unmarshal Head message
	r.Header.Load(header)

	if r.Header.MessageSize <= 0 {
		// maybe heart beat data message, so do ignore here
		return ERR_IGNORE_ERR
	}

	// read left
	bodySize := r.Header.MessageSize - r.Header.MetaSize
	meta := make([]byte, r.Header.MetaSize)
	body := make([]byte, bodySize)

	if l, err := io.ReadAtLeast(rw, meta, int(r.Header.MetaSize)); err != nil {
		return fmt.Errorf("Read incomplete meta, expect %d, read %d", r.Header.MetaSize, l)
	}

	n, err := io.ReadFull(rw, body)
	if err != nil {
		return err
	}
	if n != int(bodySize) {
		return fmt.Errorf("Read incomplete message, expect %d, read %d", r.Header.MessageSize, n)
	}

	proto.Unmarshal(meta, &r.Meta)

	attachmentSize := r.Meta.GetAttachmentSize()
	dataSize := bodySize - attachmentSize

	if dataSize > 0 {
		r.Data = body[0:dataSize]

		switch r.Meta.GetCompressType() {
		case COMPRESS_GZIP:
			r.Data, err = GUNZIP(r.Data)
			if err != nil {
				return err
			}
		case COMPRESS_SNAPPY:
			dst := make([]byte, 1)
			r.Data, err = snappy.Decode(dst, r.Data)
			if err != nil {
				return err
			}
		}
	}
	// if need read Attachment
	if attachmentSize > 0 {
		r.Attachment = body[bodySize-attachmentSize : bodySize]
	}

	return nil
}

// Load 加载[]byte里的数据至Package
func (r *Package) Load(b []byte) error {
	if b == nil {
		return errors.New("b is nil")
	}

	buf := bytes.NewBuffer(b)

	return r.ReadIO(buf)
}
