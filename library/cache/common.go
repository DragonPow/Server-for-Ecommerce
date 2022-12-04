package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
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

func UnMarshal[T any](value []byte, result any) error {
	reader := bytes.Buffer{}
	g, err := gzip.NewReader(&reader)
	if err != nil {
		return err
	}
	defer g.Close()

	_, err = g.Read(value)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(reader.Bytes(), &result); err != nil {
		return err
	}
	return nil
}
