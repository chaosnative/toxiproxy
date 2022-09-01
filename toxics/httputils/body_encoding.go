package httputils

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
)

func Gzip(body []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(body); err != nil {
		gz.Close()
		return nil, err
	}
	gz.Close()
	return buf.Bytes(), nil
}

// Deflate compresses the body using the DEFLATE algorithm.
func Deflate(body []byte) ([]byte, error) {
	var buf bytes.Buffer
	z := zlib.NewWriter(&buf)
	if _, err := z.Write(body); err != nil {
		z.Close()
		return nil, err
	}
	z.Close()
	return buf.Bytes(), nil
}
