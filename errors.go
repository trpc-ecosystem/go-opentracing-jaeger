//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Package jaeger is the trpc-opentracing-jaeger plugin.
package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"trpc.group/trpc-go/trpc-go/errs"
)

// setErrorTags sets one or more tags on the given span according to the
// error
func setErrorTags(span opentracing.Span, err error) {
	code := errs.Code(err)
	span.SetTag("response.code", code)
	ext.Error.Set(span, true)
}
