package util

import (
	"encoding/json"
	"io"
)

func Decode[T interface{}](reader io.Reader, val T) error {
	return json.NewDecoder(reader).Decode(val)
}

func Encode[T interface{}](writer io.Writer, val T) error {
	return json.NewEncoder(writer).Encode(val)
}
