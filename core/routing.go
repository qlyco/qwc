package core

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

type pageRoute struct {
	file       string
	source     string
	url        string
	pageNumber int
	isText     bool
}

func (server *Server) findRoute(urlPath string) pageRoute {
	route := server.findSource(urlPath)

	if route.file != "" {
		return route
	}

	baseURL, pageNumber, isNumbered := splitPageNumber(urlPath)

	if !isNumbered {
		return pageRoute{}
	}

	route = server.findSource(baseURL)
	route.pageNumber = pageNumber

	return route
}

func (server *Server) findSource(urlPath string) pageRoute {
	exact := server.newRoute(strings.TrimPrefix(urlPath, res.PathSeparator))

	if exact.isText && isFile(exact.file) {
		return exact
	}

	names := []string{path.Join(urlPath, res.IndexName), urlPath}

	if path.Ext(urlPath) == res.PageExt {
		names = []string{strings.TrimSuffix(urlPath, res.PageExt)}
	}

	for _, name := range names {
		route := server.findPage(name)

		if route.file != "" && !route.isText && !route.isSpecial() {
			return route
		}
	}

	return pageRoute{}
}

func (server *Server) findPage(name string) pageRoute {
	for _, extension := range res.PageExts {
		route := server.newRoute(strings.TrimPrefix(fmt.Sprintf(res.JoinFormat, name, extension), res.PathSeparator))

		if isFile(route.file) {
			return route
		}
	}

	return pageRoute{}
}

func (server *Server) walkPages(visit func(route pageRoute) error) error {
	return walkPaths(server.pages, func(relativePath string, filePath string) error {
		route := server.newRoute(relativePath)

		if !route.isText && server.findPage(trimExt(relativePath)).source != relativePath {
			return nil
		}

		return visit(route)
	})
}

func (server *Server) newRoute(source string) pageRoute {
	route := pageRoute{
		file:   filepath.Join(server.pages, filepath.FromSlash(source)),
		source: source,
		url:    fmt.Sprintf(res.RootedPathFormat, source),
	}

	route.isText = isTextPage(source) || (isPageFile(source) && hasTextType(route.file))

	if route.isText {
		return route
	}

	route.url = fmt.Sprintf(res.RootedPathFormat, trimExt(source))

	if path.Base(trimExt(source)) == res.IndexName {
		route.url = strings.TrimSuffix(route.url, res.IndexName)
	}

	return route
}

func (route pageRoute) sitePath() string {
	if route.pageNumber == 0 && route.isText {
		return route.source
	}

	if route.pageNumber == 0 {
		return fmt.Sprintf(res.JoinFormat, trimExt(route.source), res.PageExt)
	}

	return path.Join(strings.TrimPrefix(route.url, res.PathSeparator), res.PageSegment, strconv.Itoa(route.pageNumber), res.IndexPage)
}

func splitPageNumber(urlPath string) (string, int, bool) {
	parent, last := path.Split(urlPath)
	baseURL, segment := path.Split(path.Clean(parent))
	pageNumber, err := strconv.Atoi(last)

	if err != nil || segment != res.PageSegment || pageNumber < 1 || strconv.Itoa(pageNumber) != last {
		return "", 0, false
	}

	return path.Clean(baseURL), pageNumber, true
}

func pageNumberURL(baseURL string, pageNumber int) string {
	if pageNumber <= 1 {
		return baseURL
	}

	return fmt.Sprintf(res.PageNumberURLFormat, strings.TrimSuffix(baseURL, res.PathSeparator), res.PageSegment, pageNumber)
}

func (route pageRoute) isSpecial() bool {
	name := path.Base(trimExt(route.source))

	return !route.isText && (name == res.NotFoundName || name == res.FallbackName)
}

func isTextPage(source string) bool {
	_, isText := textType(source)

	return isText && !isPageFile(source)
}

func isPageFile(source string) bool {
	return slices.Contains(res.PageExts, strings.ToLower(path.Ext(source)))
}

func hasTextType(filePath string) bool {
	values, err := readFrontMatter(filePath)

	return err == nil && values.isText()
}

func trimExt(name string) string {
	return strings.TrimSuffix(name, path.Ext(name))
}
