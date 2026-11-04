package core

import (
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func (values templateData) parseTags() {
	text, isString := values.text(res.TagsKey)

	if !isString {
		return
	}

	tags := []string{}

	for _, name := range strings.Split(text, res.ListSeparator) {
		name = strings.TrimSpace(name)

		if !isUsableTag(name) || slices.Contains(tags, name) {
			continue
		}

		tags = append(tags, name)
	}

	values[res.TagsKey] = tags
}

func isUsableTag(name string) bool {
	return name != "" && !strings.Contains(name, res.PathSeparator) && name != res.CurrentDir && name != res.ParentDir
}

func entryTags(values templateData) []string {
	tags, _ := values[res.TagsKey].([]string)

	return tags
}

func (values templateData) normalizeTagsSlug() {
	slug, _ := values.text(res.TagsSlugKey)
	slug = strings.Trim(strings.TrimSpace(slug), res.PathSeparator)

	if !isUsableTag(slug) {
		slug = res.TagsDir
	}

	values[res.TagsSlugKey] = slug
}

func (server *Server) tagBase() string {
	sitewide, err := server.loadSitewideData()

	if err != nil {
		return res.TagsDir
	}

	slug, _ := sitewide.text(res.TagsSlugKey)

	if !isUsableTag(slug) {
		return res.TagsDir
	}

	return slug
}

func (server *Server) tagged(name string) ([]templateData, error) {
	entries, err := server.listPages(res.AnySegments, false)

	if err != nil {
		return nil, err
	}

	matches := []templateData{}

	for _, entry := range entries {
		if slices.Contains(entryTags(entry), name) {
			matches = append(matches, entry)
		}
	}

	return matches, nil
}

func (server *Server) tags() ([]templateData, error) {
	entries, err := server.listPages(res.AnySegments, false)

	if err != nil {
		return nil, err
	}

	counts := map[string]int{}

	for _, entry := range entries {
		for _, name := range entryTags(entry) {
			counts[name]++
		}
	}

	base := server.tagBase()
	names := slices.Sorted(maps.Keys(counts))
	infos := make([]templateData, 0, len(names))

	for _, name := range names {
		infos = append(infos, templateData{
			"Name":  name,
			"Count": counts[name],
			"URL":   fmt.Sprintf(res.TagURLFormat, base, name),
		})
	}

	return infos, nil
}

func (server *Server) findTagRoute(urlPath string) pageRoute {
	urlPath = path.Clean(urlPath)
	base := res.PathSeparator + server.tagBase()

	if urlPath == base {
		if !server.hasLayout(res.TaxonomyLayout) {
			return pageRoute{}
		}

		return server.taxonomyRoute()
	}

	name, isTagPath := strings.CutPrefix(urlPath, base+res.PathSeparator)

	if !isTagPath || !server.hasLayout(res.TermLayout) {
		return pageRoute{}
	}

	tags, err := server.tags()

	if err != nil {
		server.logger.Printf(res.FileLogFormat, server.sitePath(server.pages), err)
		return pageRoute{}
	}

	for _, tag := range tags {
		tagName, _ := tag.text("Name")

		if tagName == name {
			return server.termRoute(name)
		}
	}

	return pageRoute{}
}

func (server *Server) walkTagRoutes(visit func(route pageRoute) error) error {
	routes, err := server.tagRoutes()

	if err != nil {
		return err
	}

	for _, route := range routes {
		err = visit(route)

		if err != nil {
			return err
		}
	}

	return nil
}

func (server *Server) tagRoutes() ([]pageRoute, error) {
	tags, err := server.tags()

	if err != nil {
		return nil, err
	}

	var routes []pageRoute

	if server.hasLayout(res.TaxonomyLayout) {
		routes = append(routes, server.taxonomyRoute())
	}

	if server.hasLayout(res.TermLayout) {
		for _, tag := range tags {
			tagName, _ := tag.text("Name")
			routes = append(routes, server.termRoute(tagName))
		}
	}

	return slices.DeleteFunc(routes, func(route pageRoute) bool {
		return server.findSource(path.Clean(route.url)).file != ""
	}), nil
}

func (server *Server) taxonomyRoute() pageRoute {
	source := path.Join(server.tagBase(), res.IndexPage)

	return server.generatedRoute(source, res.TaxonomyLayout, templateData{}, server.tags)
}

func (server *Server) termRoute(tag string) pageRoute {
	source := path.Join(server.tagBase(), tag, res.IndexPage)

	return server.generatedRoute(source, res.TermLayout, templateData{res.TagKey: tag}, func() ([]templateData, error) {
		return server.tagged(tag)
	})
}

func (server *Server) generatedRoute(source string, layout string, values templateData, list func() ([]templateData, error)) pageRoute {
	route := server.newRoute(source)
	layoutValues := server.layoutValues(layout)

	values[res.LayoutKey] = layout
	values[res.URLKey] = route.url

	pageSize, isString := layoutValues.text(res.PageSizeKey)

	if isString {
		values[res.PageSizeKey] = pageSize
	}

	route.values = values
	route.list = list

	return route
}

func (server *Server) layoutValues(layout string) templateData {
	values, err := readFrontMatter(filepath.Join(server.layouts, filepath.FromSlash(fmt.Sprintf(res.JoinFormat, layout, res.PageExt))))

	if err != nil {
		return templateData{}
	}

	return values
}

func (server *Server) hasLayout(name string) bool {
	return isFile(filepath.Join(server.layouts, filepath.FromSlash(fmt.Sprintf(res.JoinFormat, name, res.PageExt))))
}

func (server *Server) loadRouteValues(route pageRoute) (templateData, []byte, error) {
	if route.values != nil {
		return maps.Clone(route.values), nil, nil
	}

	return loadPage(route.file, pageNumberURL(route.url, route.pageNumber))
}
