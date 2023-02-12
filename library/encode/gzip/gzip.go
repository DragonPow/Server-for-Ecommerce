package gzip

import (
	"Server-for-Ecommerce/library/encode"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
)

type myGzip struct {
}

func NewGzipEncoder() encode.HttpEncoder {
	return &myGzip{}
}

func (g *myGzip) Marshal(value any) ([]byte, error) {
	buf := bytes.Buffer{}
	w := gzip.NewWriter(&buf)
	defer w.Close()

	b, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *myGzip) Unmarshal(value []byte, result any) error {
	reader := bytes.NewBuffer(value)
	r, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return err
	}
	return nil
}
