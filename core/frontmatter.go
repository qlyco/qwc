package core

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/qlyco/qwc/internal/res"
)

var dateLayouts = []string{time.DateOnly, time.RFC3339}

func loadPage(filePath string, url string) (templateData, []byte, error) {
	source, err := os.ReadFile(filePath)

	if err != nil {
		return nil, nil, err
	}

	values, body := parseFrontMatter(source)
	values[res.URLKey] = url

	return values, body, nil
}

func readFrontMatter(filePath string) (templateData, error) {
	handle, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer handle.Close()

	var source []byte

	chunk := make([]byte, res.FrontMatterChunkSize)

	for !isFrontMatterComplete(source) {
		count, err := handle.Read(chunk)
		source = append(source, chunk[:count]...)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	values, _ := parseFrontMatter(source)

	return values, nil
}

func isFrontMatterComplete(source []byte) bool {
	text := bytes.TrimLeft(source, res.LeadingSpace)

	if len(text) < len(res.FrontMatterStart) {
		return false
	}

	if !bytes.HasPrefix(text, []byte(res.FrontMatterStart)) {
		return true
	}

	return bytes.Contains(text, []byte(res.FrontMatterEnd))
}

func parseFrontMatter(source []byte) (templateData, []byte) {
	text := strings.TrimLeft(string(source), res.LeadingSpace)

	if !strings.HasPrefix(text, res.FrontMatterStart) {
		return templateData{}, source
	}

	header, body, isClosed := strings.Cut(strings.TrimPrefix(text, res.FrontMatterStart), res.FrontMatterEnd)

	if !isClosed {
		return templateData{}, source
	}

	values := parseKeyValues(header)
	values.parseDate()

	return values, []byte(strings.TrimLeft(body, res.LeadingSpace))
}

func parseKeyValues(header string) templateData {
	values := templateData{}

	for _, line := range strings.Split(header, res.LineBreak) {
		key, value, isPair := strings.Cut(line, res.KeyValueSeparator)

		if isPair {
			values[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}

	return values
}

func (values templateData) parseDate() {
	text, isString := values.text(res.DateKey)

	if !isString {
		return
	}

	for _, layout := range dateLayouts {
		date, err := time.Parse(layout, text)

		if err == nil {
			values[res.DateKey] = date
			return
		}
	}
}
