//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 Tencent.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Package jaeger trpc-opentracing-jaeger plugin.
package jaeger

import (
	"context"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/plugin"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracelog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/config"
)

const (
	pluginName = "jaeger"
	pluginType = "tracing"
)

var (
	tRPCComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "tRPC"}
)

func init() {
	plugin.Register(pluginName, &jaegerPlugin{})
}

// jaegerPlugin jaeger trpc plugin implementation.
type jaegerPlugin struct{}

// PluginType jaeger trpc plugin name.
func (p *jaegerPlugin) Type() string {
	return pluginType
}

// Setup jaeger instance initialization.
func (p *jaegerPlugin) Setup(name string, decoder plugin.Decoder) error {
	cfg := config.Configuration{}
	err := decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	filter.Register("jaeger", ServerFilter(tracer), ClientFilter(tracer))
	return nil
}

type metadataTextMap codec.MetaData

// Set implements the opentracing.TextMapReader interface.
func (m metadataTextMap) Set(key, val string) {
	m[key] = []byte(val)
}

// ForeachKey implements the opentracing.TextMapReader interface.
func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, v := range m {
		if err := callback(k, string(v)); err != nil {
			return err
		}
	}
	return nil
}

// ServerFilter is a distributed trace filter for server-side RPC.
func ServerFilter(tracer opentracing.Tracer) filter.ServerFilter {
	return func(ctx context.Context, req interface{}, handler filter.ServerHandleFunc) (interface{}, error) {
		msg := codec.Message(ctx)
		md := msg.ServerMetaData()
		parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Errorf("trpc-opentracing-jaeger: failed parsing trace information: %v", err)
		}
		log.Trace("trpc-opentracing-jaeger: server metadata ", md)
		serverSpan := tracer.StartSpan(
			msg.ServerRPCName(),
			ext.RPCServerOption(parentSpanContext),
			tRPCComponentTag,
		)

		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		rsp, err := handler(ctx, req)
		if err != nil {
			setErrorTags(serverSpan, err)
			serverSpan.LogFields(tracelog.String("event", "error"), tracelog.String("message", err.Error()))
		}
		serverSpan.Finish()
		return rsp, err
	}
}

// ClientFilter is a distributed trace filter for client-side RPC.
func ClientFilter(tracer opentracing.Tracer) filter.ClientFilter {
	return func(ctx context.Context, req, rsp interface{}, handler filter.ClientHandleFunc) error {
		var parentSpanCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentSpanCtx = parent.Context()
		}

		opts := []opentracing.StartSpanOption{
			opentracing.ChildOf(parentSpanCtx),
			ext.SpanKindRPCClient,
			tRPCComponentTag,
		}
		msg := codec.Message(ctx)
		clientSpan := tracer.StartSpan(msg.ClientRPCName(), opts...)

		md := msg.ClientMetaData().Clone()
		if len(md) == 0 {
			md = codec.MetaData{}
		}
		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, metadataTextMap(md)); err != nil {
			log.Errorf("trpc-opentracing-jaeger: failed sericalizing trace information: %v", err)
		}
		log.Trace("trpc-opentracing-jaeger: client metadata ", md)
		msg.WithClientMetaData(md)
		ctx = opentracing.ContextWithSpan(ctx, clientSpan)
		err := handler(ctx, req, rsp)
		if err != nil {
			setErrorTags(clientSpan, err)
			clientSpan.LogFields(tracelog.String("event", "error"), tracelog.String("message", err.Error()))
		}
		clientSpan.Finish()
		return err
	}
}
