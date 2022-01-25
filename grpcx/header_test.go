package pbrpc_test

import (
	"bytes"
	"fmt"
	"github.com/layasugar/laya/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

func TestRpcDataWriteReader(t *testing.T) {

	h := grpc.Header{}
	h.SetMagicCode([]byte("PRPB"))
	h.MessageSize = 12300
	h.MetaSize = 59487

	bs := h.Bytes()

	if len(bs) != pbrpc.HeaderSize {
		t.Errorf("current head size is '%d', should be '%d'", len(bs), pbrpc.HeaderSize)
	}

	h2 := pbrpc.Header{}
	h2.Load(bs)
	if !bytes.Equal(h.MagicCode, h2.MagicCode) {
		t.Errorf("magic code is not same. expect '%b' actual is '%b'", h.MagicCode, h2.MagicCode)
	}

	assert.Equal(t, h.MessageSize, h2.MessageSize, fmt.Sprintf("expect message size is %d, acutal value is %d", h.MessageSize, h2.MessageSize))

}
