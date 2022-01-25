package grpcx_test

import (
	"bytes"
	"fmt"
	"github.com/layasugar/laya/grpcx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRpcDataWriteReader(t *testing.T) {
	h := grpcx.Header{}
	h.SetMagicCode([]byte("PRPB"))
	h.MessageSize = 12300
	h.MetaSize = 59487

	bs := h.Bytes()

	if len(bs) != grpcx.HeaderSize {
		t.Errorf("current head size is '%d', should be '%d'", len(bs), grpcx.HeaderSize)
	}

	h2 := grpcx.Header{}
	h2.Load(bs)
	if !bytes.Equal(h.MagicCode, h2.MagicCode) {
		t.Errorf("magic code is not same. expect '%b' actual is '%b'", h.MagicCode, h2.MagicCode)
	}

	assert.Equal(t, h.MessageSize, h2.MessageSize, fmt.Sprintf("expect message size is %d, acutal value is %d", h.MessageSize, h2.MessageSize))
}
