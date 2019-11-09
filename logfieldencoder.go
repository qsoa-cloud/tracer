package tracer

import (
	"encoding/json"
	log2 "log"
	"strconv"

	"github.com/opentracing/opentracing-go/log"
)

type LogFieldEncoder struct {
	K string `json:"k"`
	T string `json:"t"`
	V string `json:"v"`
}

func (e *LogFieldEncoder) EmitString(key, value string) {
	e.K = key
	e.T = "string"
	e.V = value
}

func (e *LogFieldEncoder) EmitBool(key string, value bool) {
	e.K = key
	e.T = "bool"
	if value {
		e.V = "t"
	} else {
		e.V = "f"
	}
}

func (e *LogFieldEncoder) EmitInt(key string, value int) {
	e.K = key
	e.T = "int"
	e.V = strconv.FormatInt(int64(value), 10)
}

func (e *LogFieldEncoder) EmitInt32(key string, value int32) {
	e.K = key
	e.T = "int32"
	e.V = strconv.FormatInt(int64(value), 10)
}

func (e *LogFieldEncoder) EmitInt64(key string, value int64) {
	e.K = key
	e.T = "int64"
	e.V = strconv.FormatInt(value, 10)
}

func (e *LogFieldEncoder) EmitUint32(key string, value uint32) {
	e.K = key
	e.T = "uint32"
	e.V = strconv.FormatUint(uint64(value), 10)
}

func (e *LogFieldEncoder) EmitUint64(key string, value uint64) {
	e.K = key
	e.T = "uint64"
	e.V = strconv.FormatUint(value, 10)
}

func (e *LogFieldEncoder) EmitFloat32(key string, value float32) {
	e.K = key
	e.T = "float32"
	e.V = strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func (e *LogFieldEncoder) EmitFloat64(key string, value float64) {
	e.K = key
	e.T = "float64"
	e.V = strconv.FormatFloat(value, 'f', -1, 64)
}

func (e *LogFieldEncoder) EmitObject(key string, value interface{}) {
	e.K = key
	e.T = "object"

	j, err := json.Marshal(value)
	if err != nil {
		log2.Printf("Cannot emit object: %s", err.Error())
	}
	e.V = string(j)
}

func (e *LogFieldEncoder) EmitLazyLogger(value log.LazyLogger) {
	value(e)
}
