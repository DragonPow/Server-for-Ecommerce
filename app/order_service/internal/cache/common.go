package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
)

func Marshal[T any](value T) ([]byte, error) {
	buf := bytes.Buffer{}
	g := gzip.NewWriter(&buf)
	defer g.Close()

	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	_, err = g.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func UnMarshal[T any](value interface{}) (T, error) {
	reader := bytes.Buffer{}
	g, err := gzip.NewReader(&reader)
	if err != nil {
		return nil, err
	}
	defer g.Close()

	b, ok := value.([]byte)
	if !ok {
		return nil, errors.New("value must be type []byte")
	}
	_, err = g.Read(b)
	if err != nil {
		return nil, err
	}
	var result T
	if err = json.Unmarshal(reader.Bytes(), &result); err != nil {
		return nil, err
	}
	return result, nil
}
