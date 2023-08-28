// Package jaeger trpc-opentracing-jaeger 插件
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
