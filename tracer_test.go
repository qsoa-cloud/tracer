package tracer_test

import (
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/harness"

	"gopkg.qsoa.cloud/tracer"
)

type devNull struct{}

func (devNull) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestTracer_ApiChecks(t *testing.T) {
	harness.RunAPIChecks(t, func() (opentracing.Tracer, func()) {
		return tracer.New(devNull{}), nil
	}, harness.CheckEverything(), harness.CheckBaggageValues(false))
}
