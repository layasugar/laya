package pbrpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	COMPRESS_NO     int32 = 0
	COMPRESS_SNAPPY int32 = 1
	COMPRESS_GZIP   int32 = 2

	HeaderSize = 12
	MagicCode  = "PRPC"
)

// Header PbRPC header content
type Header struct {
	MagicCode   []byte
	MessageSize int32
	MetaSize    int32
}

// NewHeader create new empty header
func NewHeader() *Header {
	return &Header{
		MagicCode:   []byte(MagicCode),
		MessageSize: 0,
		MetaSize:    0,
	}
}

// Bytes return header struct to byte array
func (h *Header) Bytes() []byte {
	b := new(bytes.Buffer)

	_ = binary.Write(b, binary.BigEndian, h.MagicCode)
	_ = binary.Write(b, binary.BigEndian, intToBytes(h.MessageSize))
	_ = binary.Write(b, binary.BigEndian, intToBytes(h.MetaSize))
	return b.Bytes()
}

// Load use bytes array to update field for header
func (h *Header) Load(data []byte) error {
	if data == nil || len(data) != HeaderSize {
		return fmt.Errorf("data should be a %d-bytes slice", HeaderSize)
	}
	h.MagicCode = data[0:4]
	h.MessageSize = int32(binary.BigEndian.Uint32(data[4:8]))
	h.MetaSize = int32(binary.BigEndian.Uint32(data[8:12]))
	return nil
}

// SetMagicCode change magic code in header
func (h *Header) SetMagicCode(MagicCode []byte) error {
	if MagicCode == nil || len(MagicCode) != 4 {
		return fmt.Errorf("MagicCode should be a 4-bytes slice")
	}
	h.MagicCode = MagicCode
	return nil
}

func intToBytes(i int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))
	return b
}
