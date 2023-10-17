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

package jaeger

import (
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"trpc.group/trpc-go/trpc-go/errs"
	trpcpb "trpc.group/trpc/trpc-protocol/pb/go/trpc"
)

// TestSpanTags
func TestSpanTags(t *testing.T) {
	tracer := mocktracer.New()
	retCodes := []trpcpb.TrpcRetCode{
		errs.RetServerDecodeFail,
		errs.RetClientDecodeFail,
		errs.RetClientNetErr,
		errs.RetServerNoFunc,
		errs.RetServerSystemErr,
		errs.RetServerTimeout,
	}

	for _, retCode := range retCodes {
		tracer.Reset()
		span := tracer.StartSpan("test-trace-client")
		err := errs.New(retCode, "")
		setErrorTags(span, err)
		span.Finish()

		// Assert added tags
		rawSpan := tracer.FinishedSpans()[0]
		expectedTags := map[string]interface{}{
			"response.code": retCode,
		}
		if err != nil {
			expectedTags["error"] = true
		}
		assert.Equal(t, expectedTags, rawSpan.Tags())

		// Server error
		tracer.Reset()
		span = tracer.StartSpan("test-trace-server")
		err = errs.New(retCode, "")
		setErrorTags(span, err)
		span.Finish()

		// Assert added tags
		rawSpan = tracer.FinishedSpans()[0]
		expectedTags = map[string]interface{}{
			"response.code": retCode,
		}
		if err != nil {
			expectedTags["error"] = true
		}
		assert.Equal(t, expectedTags, rawSpan.Tags())
	}
}
