# Golang Jaeger Example

Due to a lack of examples on how to connect redis and gin together and add a trace id to response headers, I decided to create this small example.

## Requirements

- VSCode
- VSCode remote container plugin

## Installation

1. Run remote container with the .devcontainer config

2. cURL <http://localhost:29090>

3. Open Jaeger UI at <http://localhost:16686>

4. See a single trace with 2 spans under it, one for Gin and one for Redis

## Step By Step Code Guide

- Add global tracer with jaeger config
- Add go jaeger gin middleware
- Add redis/v8 hook for tracing
- Extract span on each incoming request
- Pass context with span to redis / any other service you need
- Extract Trace Id
- Pass Trace Id as response header for client
