package tracer

import "github.com/opentracing/opentracing-go"

var _ opentracing.SpanContext = &SpanContext{}

type SpanContext struct {
	TraceID      uint64            `json:"t"`
	SpanID       uint64            `json:"s"`
	ParentSpanId uint64            `json:"p,omitempty"`
	Baggage      map[string]string `json:"-"`
}

func (c *SpanContext) ForeachBaggageItem(handler func(k, v string) bool) {
	for k, v := range c.Baggage {
		if !handler(k, v) {
			break
		}
	}
}
