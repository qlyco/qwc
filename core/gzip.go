package core

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func compress(body []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)

	if err != nil {
		return nil, err
	}

	_, err = writer.Write(body)
	err = errors.Join(err, writer.Close())

	return buffer.Bytes(), err
}

func acceptsGzip(request *http.Request) bool {
	for _, encoding := range strings.Split(request.Header.Get(res.AcceptEncodingHeader), res.ListSeparator) {
		name, parameters, _ := strings.Cut(encoding, res.ParameterSeparator)

		if strings.TrimSpace(name) == res.GzipEncoding && strings.TrimSpace(parameters) != res.RefusedQuality {
			return true
		}
	}

	return false
}

func decompress(body []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	return io.ReadAll(reader)
}
