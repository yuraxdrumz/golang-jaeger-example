package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

func OpenTracing() gin.HandlerFunc {
 return func(c *gin.Context) {
	wireCtx, _ := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header))

	serverSpan := opentracing.StartSpan(c.Request.URL.Path, ext.RPCServerOption(wireCtx))
	defer serverSpan.Finish()

	c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), serverSpan))

	if sc, ok := serverSpan.Context().(jaeger.SpanContext); ok {
		traceId := sc.TraceID()
		c.Header("X-Trace-Id", traceId.String())
	}
    
	c.Next()

	if c.Writer.Status() > 299 {
		ext.Error.Set(serverSpan, true)
	}

 }
}