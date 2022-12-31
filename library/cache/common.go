package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
)

func Marshal(value any) ([]byte, error) {
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

func Unmarshal(value []byte, result any) error {
	reader := bytes.NewBuffer(value)
	g, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer g.Close()

	data, err := ioutil.ReadAll(g)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return err
	}
	return nil
}
