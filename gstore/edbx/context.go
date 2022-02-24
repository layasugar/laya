// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package edbx

import (
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/layasugar/laya/core/tracex"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	tSpanName = "elasticsearch"
)

// Transport for tracing Elastic operations.
type Transport struct {
	rt http.RoundTripper
}

// Option signature for specifying options, e.g. WithRoundTripper.
type Option func(t *Transport)

// WithRoundTripper specifies the http.RoundTripper to call
// next after this transport. If it is nil (default), the
// transport will use http.DefaultTransport.
func WithRoundTripper(rt http.RoundTripper) Option {
	return func(t *Transport) {
		t.rt = rt
	}
}

// NewTransport specifies a transport that will trace Elastic
// and report back via OpenTracing.
func NewTransport(opts ...Option) *Transport {
	t := &Transport{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// RoundTrip captures the request and starts an OpenTracing span
// for Elastic PerformRequest operation.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var ctx = req.Context()
	var span opentracing.Span
	var traceCtx *tracex.TraceContext
	switch ctx.(type) {
	case *httpx.WebContext:
		if webCtx, okInterface := ctx.(*httpx.WebContext); okInterface {
			traceCtx = webCtx.TraceContext
		}
	case *grpcx.GrpcContext:
		if grpcCtx, okInterface := ctx.(*grpcx.GrpcContext); okInterface {
			traceCtx = grpcCtx.TraceContext
		}
	case *appx.Context:
		if defaultCtx, okInterface := ctx.(*appx.Context); okInterface {
			traceCtx = defaultCtx.TraceContext
		}
	}

	if traceCtx != nil {
		span = traceCtx.SpanStart(tSpanName)
		if nil != span {
			ext.Component.Set(span, "github.com/olivere/elastic/v6")
			ext.HTTPUrl.Set(span, req.URL.String())
			ext.HTTPMethod.Set(span, req.Method)
			ext.PeerHostname.Set(span, req.URL.Hostname())
			ext.PeerPort.Set(span, atouint16(req.URL.Port()))
			defer span.Finish()
		}
	}

	var (
		resp *http.Response
		err  error
	)
	if t.rt != nil {
		resp, err = t.rt.RoundTrip(req)
	} else {
		resp, err = http.DefaultTransport.RoundTrip(req)
	}
	if err != nil {
		ext.Error.Set(span, true)
	}
	if resp != nil {
		ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	}

	return resp, err
}
