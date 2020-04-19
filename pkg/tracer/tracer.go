package tracer

import (
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

// NewTracer starts and returns a opentracing tracer
func NewTracer(conf config.Tracer) opentracing.Tracer {
	switch conf.Type {
	case "jaeger":
		return StartJaegerTracer(conf.Jaeger)
	case "datadog":
		// environment variable DD_AGENT_HOST should be set to use on GCP environment
		if os.Getenv("DD_AGENT_HOST") == "" {
			return NoopTracer()
		}
		return StartDatadogTracer(conf.Datadog, os.Getenv("DD_AGENT_HOST"))
	}
	return NoopTracer()
}

// NoopTracer returns empty tracer
func NoopTracer() opentracing.Tracer {
	return opentracing.NoopTracer{}
}

// NoopSpan returns opentracing.NoopTracer{}
func NoopSpan(name string) opentracing.Span {
	return NoopTracer().StartSpan(name)
}
