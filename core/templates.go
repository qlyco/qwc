package core

import (
	"bytes"
	"fmt"
	html "html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	text "text/template"
	"text/template/parse"

	"github.com/qlyco/qwc/internal/res"
)

type templateSet[T any] interface {
	New(name string) T
	Parse(source string) (T, error)
	Clone() (T, error)
	AddParseTree(name string, tree *parse.Tree) (T, error)
	ExecuteTemplate(writer io.Writer, name string, data any) error
}

type layoutFile struct {
	name string
	body []byte
}

type pageRender struct {
	server *Server
	page   string
	body   []byte
	values templateData
	isText bool
}

func (server *Server) render(page string, body []byte, values templateData) []byte {
	functions, err := server.templateFuncs()
	job := pageRender{server: server, page: server.sitePath(page), body: body, values: values, isText: isTextPage(page) || values.isText()}

	if err != nil {
		return job.failure(nil, err)
	}

	if job.isText {
		return renderWith(job, text.New("").Funcs(functions))
	}

	return renderWith(job, html.New("").Funcs(functions))
}

func renderWith[T templateSet[T]](job pageRender, templates T) []byte {
	err := addTemplates(templates, job.server.shared)

	if err != nil {
		return job.failure(nil, err)
	}

	layouts := job.server.layoutChain(job.page, job.values)

	for _, layout := range slices.Backward(layouts) {
		_, err = templates.New(layoutName(layout.name)).Parse(string(layout.body))

		if err != nil {
			return job.failure(nil, err)
		}
	}

	content := job.body
	usesTemplates := !job.values.isDisabled(res.TemplatesKey)

	if usesTemplates {
		_, err = templates.New(job.page).Parse(string(job.body))

		if err != nil {
			return job.failure(nil, err)
		}

		content = execute(job, templates, job.page, nil)
	}

	if path.Ext(job.page) == res.MarkdownExt && !job.isText {
		content = job.convertMarkdown(content)
	}

	for level, layout := range layouts {
		if level == 0 && usesTemplates && definesContent(job.page, job.body) {
			content = execute(job, templates, layoutName(layout.name), nil)
			continue
		}

		content = execute(job, templates, layoutName(layout.name), textTree(content))
	}

	return content
}

func execute[T templateSet[T]](job pageRender, templates T, name string, content *parse.Tree) []byte {
	var output bytes.Buffer

	clone, err := templates.Clone()

	if err == nil && content != nil {
		_, err = clone.AddParseTree(res.ContentBlock, content)
	}

	if err == nil {
		err = clone.ExecuteTemplate(&output, name, job.values)
	}

	if err != nil {
		return job.failure(output.Bytes(), err)
	}

	return output.Bytes()
}

func (job pageRender) convertMarkdown(content []byte) []byte {
	output, err := convertMarkdown(content)

	if err != nil {
		return job.failure(content, err)
	}

	return output
}

func (job pageRender) failure(output []byte, err error) []byte {
	job.server.logger.Printf(res.FileLogFormat, job.page, err)

	message := res.ErrorMarker

	if job.server.verbose {
		message = fmt.Sprintf(res.ErrorMarkerFormat, err)
	}

	if !job.isText {
		message = html.HTMLEscapeString(message)
	}

	return append(slices.Clone(output), message...)
}

func (server *Server) layoutChain(page string, values templateData) []layoutFile {
	var layouts []layoutFile

	name, hasLayout := values.text(res.LayoutKey)

	for hasLayout {
		isLoop := slices.ContainsFunc(layouts, func(layout layoutFile) bool { return layout.name == name })

		if isLoop {
			server.logger.Printf(res.LayoutLoopMessage, name, page)
			return layouts
		}

		source, err := os.ReadFile(filepath.Join(server.layouts, filepath.FromSlash(fmt.Sprintf(res.JoinFormat, name, res.PageExt))))

		if err != nil {
			server.logger.Printf(res.MissingLayoutMessage, name, page)
			return layouts
		}

		layoutValues, body := parseFrontMatter(source)
		layouts = append(layouts, layoutFile{name: name, body: body})
		name, hasLayout = layoutValues.text(res.LayoutKey)
	}

	return layouts
}

func layoutName(name string) string {
	return fmt.Sprintf(res.TemplateNameFormat, res.LayoutsDir, name)
}

func addTemplates[T templateSet[T]](templates T, directories []string) error {
	for _, directory := range directories {
		prefix := filepath.Base(directory)

		err := walkFiles(directory, res.PageExt, func(name string, source []byte) error {
			_, err := templates.New(fmt.Sprintf(res.TemplateNameFormat, prefix, name)).Parse(string(source))

			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func definesContent(page string, body []byte) bool {
	trees := map[string]*parse.Tree{}
	tree := parse.New(page)
	tree.Mode = parse.SkipFuncCheck

	_, err := tree.Parse(string(body), "", "", trees)
	content, isDefined := trees[res.ContentBlock]

	return err == nil && isDefined && !parse.IsEmptyTree(content.Root)
}

func textTree(content []byte) *parse.Tree {
	node := &parse.TextNode{NodeType: parse.NodeText, Text: content}
	list := &parse.ListNode{NodeType: parse.NodeList, Nodes: []parse.Node{node}}

	return &parse.Tree{Name: res.ContentBlock, Root: list}
}

func (server *Server) templateFuncs() (text.FuncMap, error) {
	sitewide, err := server.loadSitewideData()

	if err != nil {
		return nil, err
	}

	return text.FuncMap{
		res.DictFunc:        dict,
		res.SliceFunc:       slice,
		res.SitewideFunc:    func() templateData { return sitewide },
		res.CollectionFunc:  server.collection,
		res.PagesFunc:       server.allPages,
		res.TaggedFunc:      server.tagged,
		res.TagsFunc:        server.tags,
		res.UpperCaseFunc:   upperCase,
		res.LowerCaseFunc:   lowerCase,
		res.TitleCaseFunc:   titleCase,
		res.ReplaceFunc:     replace,
		res.SplitFunc:       split,
		res.JoinFunc:        join,
		res.PluginFunc:      server.plugin,
		res.RelativeURLFunc: server.address.relativeURL,
		res.AbsoluteURLFunc: server.address.absoluteURL,
	}, nil
}
