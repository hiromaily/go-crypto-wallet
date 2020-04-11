package tracer

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

type TraceLog struct {
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
}

func StartDatadogTracer(conf config.TracerDetail, hostName string) opentracing.Tracer {
	opts := []tracer.StartOption{
		tracer.WithServiceName(conf.ServiceName),
		tracer.WithAgentAddr(fmt.Sprintf("%s%s", hostName, conf.CollectorEndpoint)),
		tracer.WithAnalyticsRate(conf.SamplingProbability),
	}
	if conf.IsDebug {
		opts = append(opts, tracer.WithDebugMode(true))
	}

	return opentracer.New(opts...)
}
