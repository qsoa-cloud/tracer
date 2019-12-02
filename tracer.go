package tracer

import (
	"encoding/json"
	"fmt"
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
	case opentracing.HTTPHeaders, opentracing.TextMap:
		carrier, ok := carrier.(opentracing.TextMapWriter)
		if !ok {
			return opentracing.ErrInvalidCarrier
		}
		carrier.Set("X-TRACE-ID", strconv.FormatUint(sCtx.TraceID, 16))
		carrier.Set("X-SPAN-ID", strconv.FormatUint(sCtx.SpanID, 16))
	case opentracing.Binary:
		carrier, ok := carrier.(io.Writer)
		if !ok {
			return opentracing.ErrInvalidCarrier
		}
		_, _ = fmt.Fprintf(carrier, "%d %d", sCtx.TraceID, sCtx.SpanID)
	default:
		return opentracing.ErrUnsupportedFormat
	}

	return nil
}

func (t *Tracer) Extract(format interface{}, opaqueCarrier interface{}) (opentracing.SpanContext, error) {
	var traceID, spanID uint64

	switch format {
	case opentracing.TextMap:

	case opentracing.HTTPHeaders:

	case opentracing.Binary:

	default:
		return nil, opentracing.ErrUnsupportedFormat
	}

	switch carrier := opaqueCarrier.(type) {
	case opentracing.TextMapCarrier:
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

	case opentracing.HTTPHeadersCarrier:
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

	case io.Reader:
		_, err := fmt.Fscanf(carrier, "%d %d", &traceID, &spanID)
		if err != nil {
			return nil, err
		}
	default:
		return nil, opentracing.ErrInvalidCarrier
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
