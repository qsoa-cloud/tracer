package tracer

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Span struct {
	tracer     *Tracer                `json:"-"`
	Operation  string                 `json:"o"`
	StartTime  time.Time              `json:"s"`
	FinishTime time.Time              `json:"f"`
	Ctx        *SpanContext           `json:"c"`
	Tags       map[string]interface{} `json:"t,omitempty"`
	JLogFields []JLogField            `json:"lf,omitempty"`
}

func NewSpan(tracer *Tracer, operation string, opts *opentracing.StartSpanOptions) *Span {
	s := &Span{
		tracer:    tracer,
		Operation: operation,
		Ctx: &SpanContext{
			SpanID: tracer.genId(),
		},
	}

	if opts == nil || opts.StartTime.IsZero() {
		s.StartTime = time.Now()
	} else {
		s.StartTime = opts.StartTime
	}

	if opts != nil {
		for _, ref := range opts.References {
			if ref.Type != opentracing.ChildOfRef {
				continue
			}
			s.Ctx.TraceID = ref.ReferencedContext.(*SpanContext).TraceID
			s.Ctx.ParentSpanId = ref.ReferencedContext.(*SpanContext).SpanID
			break
		}

		s.Tags = opts.Tags
	}

	if s.Ctx.TraceID == 0 {
		s.Ctx.TraceID = tracer.genId()
	}

	return s
}

func (s *Span) Finish() {
	s.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: time.Now(),
	})

	s.tracer.writeSpan(s)
}

func (s *Span) FinishWithOptions(opts opentracing.FinishOptions) {
	s.FinishTime = opts.FinishTime

	for _, r := range opts.LogRecords {
		s.LogFields(r.Fields...)
	}

	for _, d := range opts.BulkLogData {
		s.LogFields(d.ToLogRecord().Fields...)
	}
}

func (s *Span) Context() opentracing.SpanContext {
	return s.Ctx
}

func (s *Span) SetOperationName(operationName string) opentracing.Span {
	s.Operation = operationName

	return s
}

func (s *Span) SetTag(key string, value interface{}) opentracing.Span {
	if s.Tags == nil {
		s.Tags = map[string]interface{}{}
	}

	s.Tags[key] = value

	return s
}

func (s *Span) LogFields(fields ...log.Field) {
	jFields := make([]JLogField, len(fields))
	for i := range fields {
		jFields[i].Field = fields[i]
	}

	s.JLogFields = append(s.JLogFields, jFields...)
}

func (s *Span) LogKV(alternatingKeyValues ...interface{}) {
	fields, err := log.InterleavedKVToFields(alternatingKeyValues...)
	if err != nil {
		panic(err)
	}
	s.LogFields(fields...)
}

func (s *Span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	//panic("implement me")
	return s
}

func (s *Span) BaggageItem(restrictedKey string) string {
	//panic("implement me")
	return ""
}

func (s *Span) Tracer() opentracing.Tracer {
	return s.tracer
}

func (s *Span) LogEvent(event string) {
	s.Log(opentracing.LogData{
		Timestamp: time.Now(),
		Event:     event,
	})
}

func (s *Span) LogEventWithPayload(event string, payload interface{}) {
	s.Log(opentracing.LogData{
		Timestamp: time.Now(),
		Event:     event,
		Payload:   payload,
	})
}

func (s *Span) Log(data opentracing.LogData) {
	s.LogFields(data.ToLogRecord().Fields...)
}
