package core

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.DefinitionList, extension.Typographer),
	goldmark.WithParserOptions(parser.WithAutoHeadingID(), parser.WithAttribute()),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

func convertMarkdown(source []byte) ([]byte, error) {
	var output bytes.Buffer

	err := markdown.Convert(source, &output)

	return output.Bytes(), err
}
