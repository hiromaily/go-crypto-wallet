package tracer

import (
	"context"
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

func NewChildSpan(parentSpan *Span, name string) *Span {
	if parentSpan == nil || parentSpan.span == nil {
		return nil
	}
	ps := parentSpan.span
	span := ps.Tracer().StartSpan(
		name,
		opentracing.ChildOf(ps.Context()),
	)
	return NewSpan(span).WithLogger(parentSpan.logger)
}

func NewChildSpanFromContext(ctx context.Context, name string) *Span {
	parentSpan := opentracing.SpanFromContext(ctx)
	return NewChildSpan(NewSpan(parentSpan), name)
}

func NoopTracer() opentracing.Tracer {
	return opentracing.NoopTracer{}
}

func NoopSpan(name string) opentracing.Span {
	return NoopTracer().StartSpan(name)
}
