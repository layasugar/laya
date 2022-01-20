package grpc_test

import (
	"bytes"
	"fmt"
	"github.com/layasugar/laya/grpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRpcDataWriteReader(t *testing.T) {

	h := grpc.Header{}
	h.SetMagicCode([]byte("PRPB"))
	h.MessageSize = 12300
	h.MetaSize = 59487

	bs := h.Bytes()

	if len(bs) != grpc.HeaderSize {
		t.Errorf("current head size is '%d', should be '%d'", len(bs), grpc.HeaderSize)
	}

	h2 := grpc.Header{}
	h2.Load(bs)
	if !bytes.Equal(h.MagicCode, h2.MagicCode) {
		t.Errorf("magic code is not same. expect '%b' actual is '%b'", h.MagicCode, h2.MagicCode)
	}

	assert.Equal(t, h.MessageSize, h2.MessageSize, fmt.Sprintf("expect message size is %d, acutal value is %d", h.MessageSize, h2.MessageSize))

}
