package tracer

import (
	"context"
	"encoding/json"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	Context() context.Context
	SetError(err error)
	Finish(additionalTags ...map[string]interface{})
}

type tracu struct {
	ctx  context.Context
	span trace.Span
	tags map[string]interface{}
}

func StartTrace(ctx context.Context, opsName string) Tracer {
	tr := otel.GetTracerProvider().Tracer(opsName)
	ctx = trace.ContextWithRemoteSpanContext(
		ctx,
		trace.SpanContextFromContext(ctx),
	)

	var span trace.Span
	ctx, span = tr.Start(ctx, opsName)

	return &tracu{
		span: span,
		ctx:  ctx,
	}
}

func (t *tracu) Context() context.Context {
	return t.ctx
}

func (t *tracu) SetError(err error) {
	if err == nil {
		return
	}
	t.span.SetStatus(1, err.Error())
	t.span.RecordError(err)
}

func (t *tracu) Finish(additionalTags ...map[string]interface{}) {
	defer t.span.End()

	if additionalTags != nil && t.tags == nil {
		t.tags = make(map[string]interface{})
	}

	for _, tag := range additionalTags {
		for k, v := range tag {
			t.tags[k] = v
		}
	}

	for k, v := range t.tags {
		t.span.SetAttributes(
			attribute.String(k, toString(v)),
		)
	}
}

func toString(v interface{}) (s string) {
	switch val := v.(type) {
	case error:
		if val != nil {
			s = val.Error()
		}
	case string:
		s = val
	case int:
		s = strconv.Itoa(val)
	default:
		b, _ := json.Marshal(val)
		s = string(b)
	}

	return
}
