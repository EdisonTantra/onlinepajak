package logat

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var currentLogger AppsLogger

type AppsLogger interface {
	Info(ctx context.Context, msg string, event string, data interface{}, fields ...zap.Field)
	Debug(ctx context.Context, msg string, event string, data interface{}, fields ...zap.Field)
	Warn(ctx context.Context, msg string, event string, data interface{}, fields ...zap.Field)
	Error(ctx context.Context, msg string, event string, data interface{}, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, event string, data interface{}, fields ...zap.Field)
}

type logger struct {
	provider *zap.Logger
}

// New creates a new logger with provided options. The returned context
// contains the logger in it. Out of the box, the logger will have
// default attributes message, level, time, and caller.
func New(opts ...Option) (context.Context, AppsLogger) {
	opts = append(opts, defaultOption...)
	var opt option
	for _, f := range opts {
		f(&opt)
	}

	log := opt.getLogger()
	if log == nil {
		core := zapcore.NewCore(opt.encoderFunc(opt.config), opt.output, opt.level)
		log = &logger{provider: zap.New(core, opt.zapOption...)}
	}

	currentLogger = log
	ctx := context.WithValue(opt.ctx, contextKey, log)
	return ctx, log
}

func NewNoop() AppsLogger {
	l := zap.NewNop()
	return &logger{
		provider: l,
	}
}

func GetLogger() AppsLogger {
	return currentLogger
}

func (l *logger) Info(
	ctx context.Context, msg string, event string,
	data interface{}, fields ...zap.Field,
) {
	zapFields := []zap.Field{
		zap.String(fieldCorrelationID, correlationIDFromContext(ctx)),
		zap.Any(fieldContext, getContext(ctx)),
		//zap.String(fieldCaller, getCaller()),
		zap.String(fieldEvent, event),
		zap.Any(fieldData, data),
	}

	zapFields = append(zapFields, fields...)
	l.provider.Info(msg, zapFields...)
}

func (l *logger) Warn(
	ctx context.Context, msg string, event string,
	data interface{}, fields ...zap.Field,
) {
	zapFields := []zap.Field{
		zap.String(fieldCorrelationID, correlationIDFromContext(ctx)),
		zap.Any(fieldContext, getContext(ctx)),
		//zap.String(fieldCaller, getCaller()),
		zap.String(fieldEvent, event),
		zap.Any(fieldData, data),
	}

	zapFields = append(zapFields, fields...)
	l.provider.Warn(msg, zapFields...)
}

func (l *logger) Debug(
	ctx context.Context, msg string, event string,
	data interface{}, fields ...zap.Field,
) {
	zapFields := []zap.Field{
		zap.String(fieldCorrelationID, correlationIDFromContext(ctx)),
		zap.Any(fieldContext, getContext(ctx)),
		//zap.String(fieldCaller, getCaller()),
		zap.String(fieldEvent, event),
		zap.Any(fieldData, data),
	}

	zapFields = append(zapFields, fields...)
	l.provider.Warn(msg, zapFields...)
}

func (l *logger) Error(
	ctx context.Context, msg string, event string,
	data interface{}, fields ...zap.Field,
) {
	zapFields := []zap.Field{
		zap.String(fieldCorrelationID, correlationIDFromContext(ctx)),
		zap.Any(fieldContext, getContext(ctx)),
		//zap.String(fieldCaller, getCaller()),
		zap.String(fieldEvent, event),
		zap.Any(fieldData, data),
	}

	zapFields = append(zapFields, fields...)
	l.provider.Error(msg, zapFields...)
}

func (l *logger) Fatal(
	ctx context.Context, msg string, event string,
	data interface{}, fields ...zap.Field,
) {
	zapFields := []zap.Field{
		zap.String(fieldCorrelationID, correlationIDFromContext(ctx)),
		zap.Any(fieldContext, getContext(ctx)),
		//zap.String(fieldCaller, getCaller()),
		zap.String(fieldEvent, event),
		zap.Any(fieldData, data),
	}

	zapFields = append(zapFields, fields...)
	l.provider.Fatal(msg, zapFields...)
}

func correlationIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	sc := trace.SpanFromContext(ctx).SpanContext()
	if sc.HasTraceID() {
		return sc.TraceID().String() + "-" + sc.SpanID().String()
	}

	return ""
}
