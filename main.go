package main

import (
	"log"

	"github.com/gin-gonic/gin"
	redisopentracing "github.com/globocom/go-redis-opentracing"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/opentracing/opentracing-go"
	otgorm "github.com/smacker/opentracing-gorm"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// jaeger examples - https://opentracing.io/guides/golang/quick-start/
// jaeger redis - https://pkg.go.dev/github.com/globocom/go-redis-opentracing#readme-installation
// trace id - https://github.com/opentracing/opentracing-go/issues/188
// example of ibm for tracing - https://cloud.ibm.com/docs/go?topic=go-go-e2e-tracing

type User struct {
	gorm.Model
	Name string
}

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

    db, err := gorm.Open("sqlite3", ":memory:")
    if err != nil {
        panic(err)
    }
	db.CreateTable(&User{})
    // register callbacks must be called for a root instance of your gorm.DB
    otgorm.AddGormCallbacks(db)

	tracer := opentracing.GlobalTracer()
	hook := redisopentracing.NewHook(tracer)
	rdb.AddHook(hook)
	// create the middleware
	r := gin.Default()
	// tell gin to use the middleware
	r.Use(OpenTracing())
	r.GET("/", func(c *gin.Context) {
		_, _ = rdb.Set(c.Request.Context(), "test", 1, 0).Result()
   		db := otgorm.SetSpanToGorm(c.Request.Context(), db)
		db.Create(&User{Name: "michael"})
		c.JSON(200, "Hello world!")
	})

	r.Run(":29090")

}