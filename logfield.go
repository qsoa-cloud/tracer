package tracer

import (
	"encoding/json"

	"github.com/opentracing/opentracing-go/log"
)

type JLogField struct {
	log.Field
}

func (f JLogField) MarshalJSON() ([]byte, error) {
	s := LogFieldEncoder{}

	f.Marshal(&s)

	return json.Marshal(s)
}
