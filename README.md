# tRPC Opentracing Jaeger Plugin

[![Go Reference](https://pkg.go.dev/badge/github.com/trpc-ecosystem/go-opentracing-jaeger.svg)](https://pkg.go.dev/github.com/trpc-ecosystem/go-opentracing-jaeger)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-opentracing-jaeger)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-opentracing-jaeger)
[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://github.com/trpc-ecosystem/go-opentracing-jaeger/blob/main/LICENSE)
[![Releases](https://img.shields.io/github/release/trpc-ecosystem/go-opentracing-jaeger.svg?style=flat-square)](https://github.com/trpc-ecosystem/go-opentracing-jaeger/releases)
[![Tests](https://github.com/trpc-ecosystem/go-opentracing-jaeger/actions/workflows/prc.yml/badge.svg)](https://github.com/trpc-ecosystem/go-opentracing-jaeger/actions/workflows/prc.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-opentracing-jaeger/branch/main/graph/badge.svg)](https://app.codecov.io/gh/trpc-ecosystem/go-opentracing-jaeger/tree/main)

## Jaeger Local Installation

#### 1. Download Jaeger according to your OS

https://www.jaegertracing.io/download/

#### 2. Deploy Jaeger

Easy local deployment using jaeger-all-in-one, just start it.
```
~/Downloads/jaeger-1.13.0-darwin-amd64: ls                
example-hotrod*  jaeger-agent*  jaeger-all-in-one*  jaeger-collector*  jaeger-ingester*  jaeger-query*
~/Downloads/jaeger-1.13.0-darwin-amd64: ./jaeger-all-in-one&
```

#### 3. Business services use Jaeger

trpc_go.yaml add the following configuration:

```yaml
#Configure the server side to use jaeger filters
server:
  filter:
    - jaeger

#Configure initialized jaeger plugin configuration items
plugins:
  tracing:
    jaeger:                               
      serviceName: trpc-ecosystem
      disabled: false
      sampler:
        type: const
        param: 1
      reporter:
        localAgentHostPort: localhost:6831
```

#### 4. Client calls to business services trigger jaeger reporting

#### 5. Check out jaeger distributed tracking data in local

http://localhost:16686/

## Appendix: Jaeger Docker Installation

```
$ docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.14
```

