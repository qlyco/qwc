package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func (server *Server) loadSitewideData() (templateData, error) {
	values := templateData{}
	source, err := os.ReadFile(server.sitewide)

	if err == nil {
		err = json.Unmarshal(source, &values)
	}

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf(res.FileErrorFormat, server.sitewide, err)
	}

	data := templateData{}

	err = walkFiles(server.data, res.DataExt, func(name string, source []byte) error {
		var value any

		err := json.Unmarshal(source, &value)

		if err == nil {
			data.setNested(strings.Split(name, res.PathSeparator), value)
		}

		return err
	})

	values.normalizeTagsSlug()
	values[res.DataKey] = data

	return values, err
}

func (values templateData) setNested(keys []string, value any) {
	last := len(keys) - 1
	current := values

	for _, key := range keys[:last] {
		child, isTemplateData := current[key].(templateData)

		if !isTemplateData {
			child = templateData{}
			current[key] = child
		}

		current = child
	}

	current[keys[last]] = value
}
