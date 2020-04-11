package tracer

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Span wraps `opentracing.Span` to ensure it can be nil
type Span struct {
	span   opentracing.Span
	logger *zap.Logger
}

func NewSpan(s opentracing.Span) *Span {
	return &Span{span: s}
}

func (s *Span) WithLogger(logger *zap.Logger) *Span {
	s.logger = logger
	return s
}

func (s *Span) Tag(key string, value interface{}) {
	if s != nil {
		if s.span != nil {
			s.span.SetTag(key, value)
		}
	}
}

func (s *Span) Error(msg string, fields ...zapcore.Field) {
	if s != nil && s.logger != nil {
		s.logger.Error(msg, fields...)
	}
	s.LogSpan(msg, fields...)
}

func (s *Span) Warn(msg string, fields ...zapcore.Field) {
	if s != nil && s.logger != nil {
		s.logger.Warn(msg, fields...)
	}
	s.LogSpan(msg, fields...)
}

func (s *Span) Info(msg string, fields ...zapcore.Field) {
	if s != nil && s.logger != nil {
		s.logger.Info(msg, fields...)
	}
	s.LogSpan(msg, fields...)
}

func (s *Span) Debug(msg string, fields ...zapcore.Field) {
	if s != nil && s.logger != nil {
		s.logger.Debug(msg, fields...)
	}
	s.LogSpan(msg, fields...)
}

func (s *Span) LogSpan(msg string, fields ...zapcore.Field) {
	if s != nil && s.span != nil {
		opentracingFields := make([]opentracinglog.Field, 0, len(fields)+2)
		opentracingFields = append(opentracingFields, opentracinglog.String("event", msg))
		opentracingFields = append(opentracingFields, ZapFieldsToOpentracing(fields...)...)
		s.span.LogFields(opentracingFields...)
	}
}

func (s *Span) Finish() {
	if s != nil {
		if s.span != nil {
			s.span.Finish()
		}
	}
}

func (s *Span) NewChild(name string) *Span {
	if s != nil && s.span != nil {
		return NewChildSpan(s, name)
	}
	return nil
}

// EmptySpan is only for development
func EmptySpan() *Span {
	return NewSpan(NoopTracer().StartSpan("empty"))
}

// ZapFieldsToOpentracing returns a table of standard opentracing field based on
// the inputed table of Zap field.
func ZapFieldsToOpentracing(zapFields ...zapcore.Field) []opentracinglog.Field {
	opentracingFields := make([]opentracinglog.Field, 0, len(zapFields))

	for _, zapField := range zapFields {
		opentracingFields = append(opentracingFields, ZapFieldToOpentracing(zapField))
	}

	return opentracingFields
}

// ZapFieldToOpentracing returns a standard opentracing field based on the
// input Zap field.
func ZapFieldToOpentracing(zapField zapcore.Field) opentracinglog.Field {
	switch zapField.Type {

	case zapcore.BoolType:
		val := false
		if zapField.Integer >= 1 {
			val = true
		}
		return opentracinglog.Bool(zapField.Key, val)
	case zapcore.Float32Type:
		return opentracinglog.Float32(zapField.Key, math.Float32frombits(uint32(zapField.Integer)))
	case zapcore.Float64Type:
		return opentracinglog.Float64(zapField.Key, math.Float64frombits(uint64(zapField.Integer)))
	case zapcore.Int64Type:
		return opentracinglog.Int64(zapField.Key, int64(zapField.Integer))
	case zapcore.Int32Type:
		return opentracinglog.Int32(zapField.Key, int32(zapField.Integer))
	case zapcore.StringType:
		return opentracinglog.String(zapField.Key, zapField.String)
	case zapcore.StringerType:
		return opentracinglog.String(zapField.Key, zapField.Interface.(fmt.Stringer).String())
	case zapcore.Uint64Type:
		return opentracinglog.Uint64(zapField.Key, uint64(zapField.Integer))
	case zapcore.Uint32Type:
		return opentracinglog.Uint32(zapField.Key, uint32(zapField.Integer))
	case zapcore.DurationType:
		return opentracinglog.String(zapField.Key, time.Duration(zapField.Integer).String())
	case zapcore.ErrorType:
		return opentracinglog.Error(zapField.Interface.(error))
	case zapcore.BinaryType:
		if len(zapField.String) > 0 {
			return opentracinglog.String(zapField.Key, zapField.String)
		}
		return opentracinglog.Object(zapField.Key, zapField.Interface)
	default:
		if jmsg, ok := zapField.Interface.(json.RawMessage); ok {
			return opentracinglog.String(zapField.Key, string(jmsg))
		}
		return opentracinglog.Object(zapField.Key, zapField.Interface)
	}
}
