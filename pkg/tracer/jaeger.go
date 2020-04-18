package tracer

import (
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

// StartJaegerTracer returns jaeger tracer
func StartJaegerTracer(conf config.TracerDetail) opentracing.Tracer {
	return setJaegerTracer(conf, conf.IsDebug)
}

func setJaegerTracer(conf config.TracerDetail, isDebug bool) opentracing.Tracer {
	var suffix string
	jType := jaeger.SamplerTypeConst
	probability := conf.SamplingProbability

	if isDebug {
		suffix = "-debug"
		probability = 1
	}

	if probability < 1 && probability > 0 {
		jType = jaeger.SamplerTypeProbabilistic
	}

	cfg := jaegercfg.Configuration{
		ServiceName: conf.ServiceName + suffix,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jType,
			Param: probability, // if value is 1, it always records
		},
		Reporter: &jaegercfg.ReporterConfig{
			CollectorEndpoint: conf.CollectorEndpoint,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	opts := []jaegercfg.Option{
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	}
	if isDebug {
		opts = append(opts, jaegercfg.NoDebugFlagOnForcedSampling(false))
	}

	t, _, err := cfg.NewTracer(opts...)
	if err != nil {
		log.Fatalf("Fail to start jaeger client: %v", err)
	}

	return t
}
