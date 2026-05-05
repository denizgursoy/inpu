//go:build !goexperiment.jsonv2

package inpu

import (
	"encoding/json"
	"io"
)

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshalFromReader(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
