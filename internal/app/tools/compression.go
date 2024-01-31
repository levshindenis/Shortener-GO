package tools

import (
	"compress/gzip"
	"io"
)

func Unpacking(rbody io.ReadCloser) ([]byte, error) {
	gz, err := gzip.NewReader(rbody)
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	body, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return body, nil
}
