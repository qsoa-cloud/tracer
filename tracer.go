package tracer

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/opentracing/opentracing-go"
)

var _ opentracing.Tracer = &Tracer{}

type Tracer struct {
	writer   io.Writer
	writeMtx sync.Mutex
}

func New(w io.Writer) *Tracer {
	return &Tracer{
		writer: w,
	}
}

func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	resOpts := &opentracing.StartSpanOptions{}
	for _, opt := range opts {
		opt.Apply(resOpts)
	}

	return NewSpan(t, operationName, resOpts)
}

func (t *Tracer) Inject(sc opentracing.SpanContext, format interface{}, carrier interface{}) error {
	sCtx, ok := sc.(*SpanContext)
	if !ok {
		return opentracing.ErrInvalidSpanContext
	}

	switch format {
	case opentracing.HTTPHeaders:
		carrier, ok := carrier.(opentracing.TextMapWriter)
		if !ok {
			return opentracing.ErrInvalidCarrier
		}
		carrier.Set("X-TRACE-ID", strconv.FormatUint(sCtx.TraceID, 16))
		carrier.Set("X-SPAN-ID", strconv.FormatUint(sCtx.SpanID, 16))
	default:
		return opentracing.ErrUnsupportedFormat
	}

	return nil
}

func (t *Tracer) Extract(format interface{}, opaqueCarrier interface{}) (opentracing.SpanContext, error) {
	carrier, ok := opaqueCarrier.(opentracing.TextMapReader)
	if !ok {
		return nil, opentracing.ErrInvalidCarrier
	}

	var traceID, spanID uint64
	if err := carrier.ForeachKey(func(k, v string) error {
		var err error
		switch strings.ToLower(k) {
		case "x-trace-id":
			traceID, err = strconv.ParseUint(v, 16, 64)
			if err != nil {
				return err
			}
		case "x-span-id":
			spanID, err = strconv.ParseUint(v, 16, 64)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	}, nil
}

func (t *Tracer) genId() uint64 {
	return GetId()
}

func (t *Tracer) writeSpan(span *Span) {
	t.writeMtx.Lock()
	defer t.writeMtx.Unlock()

	if err := json.NewEncoder(t.writer).Encode(span); err != nil {
		log.Printf("encode span error: %s", err.Error())
	}
}
