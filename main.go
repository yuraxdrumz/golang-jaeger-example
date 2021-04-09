package main

import (
	"log"

	ginopentracing "github.com/Bose/go-gin-opentracing"
	"github.com/gin-gonic/gin"
	redisopentracing "github.com/globocom/go-redis-opentracing"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// jaeger examples - https://opentracing.io/guides/golang/quick-start/
// jaeger redis - https://pkg.go.dev/github.com/globocom/go-redis-opentracing#readme-installation
// trace id - https://github.com/opentracing/opentracing-go/issues/188

func main() {
	// Recommended configuration for production.
    cfg := jaegercfg.Configuration{
        ServiceName: "jaegertest",
        Sampler:     &jaegercfg.SamplerConfig{
            Type:  jaeger.SamplerTypeConst,
            Param: 1,
        },
        Reporter:    &jaegercfg.ReporterConfig{
            LogSpans: true,
			LocalAgentHostPort: "jaeger:6831",
        },
    }


	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()


	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
    })
	
	tracer := opentracing.GlobalTracer()
	hook := redisopentracing.NewHook(tracer)
	rdb.AddHook(hook)
	// create the middleware
	p := ginopentracing.OpenTracer([]byte("api-request-"))
	r := gin.Default()
	// tell gin to use the middleware
	r.Use(p)
	r.GET("/", func(c *gin.Context) {
		var span opentracing.Span
		if cspan, ok := c.Get("tracing-context"); ok {
			span = ginopentracing.StartSpanWithParent(cspan.(opentracing.Span).Context(), "helloword", c.Request.Method, c.Request.URL.Path)	
		} else {
			span = ginopentracing.StartSpanWithHeader(&c.Request.Header, "helloworld", c.Request.Method, c.Request.URL.Path)
		}

		ctx := opentracing.ContextWithSpan(c.Request.Context(), span)
		_, _ = rdb.Set(ctx, "test", 1, 0).Result()
		

		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			traceId := sc.TraceID()
			c.Header("X-Trace-Id", traceId.String())
		}
		c.JSON(200, "Hello world!")
	})

	r.Run(":29090")

}