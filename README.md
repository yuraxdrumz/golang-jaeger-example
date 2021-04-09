# Golang Jaeger Example

Due to a lack of examples on how to connect redis and gin together and add a trace id to response headers, I decided to create this small example.

## Step By Step Guide

- Add global tracer with jaeger config
- Add go jaeger gin middleware
- Add redis/v8 hook for tracing
- Extract span on each incoming request
- Pass context with span to redis / any other service you need
- Extract Trace Id
- Pass Trace Id as response header for client
