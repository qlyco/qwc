package core

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

type missingPage struct {
	route  pageRoute
	status int
}

func (server *Server) missingResponse(urlPath string, validFrom int64) pageResponse {
	page := server.findMissingPage(urlPath)

	if page.route.file == "" {
		return pageResponse{status: http.StatusNotFound}
	}

	return server.cache.text(page.route.sitePath(), page.status, validFrom, func() pageResponse {
		return server.buildMissingPage(page)
	})
}

func (server *Server) findMissingPage(urlPath string) missingPage {
	hasExtension := path.Ext(urlPath) != ""

	for directory := urlPath; ; directory = path.Dir(directory) {
		notFound := server.findPage(path.Join(directory, res.NotFoundName))

		if notFound.isSpecial() && server.isPublished(notFound) {
			return missingPage{route: notFound, status: http.StatusNotFound}
		}

		fallback := server.findPage(path.Join(directory, res.FallbackName))

		if !hasExtension && server.isActiveFallback(fallback) {
			fallback.url = strings.TrimSuffix(fallback.url, res.FallbackName)

			return missingPage{route: fallback, status: http.StatusOK}
		}

		if directory == res.PathSeparator {
			return missingPage{}
		}
	}
}

func (server *Server) isPublished(route pageRoute) bool {
	if route.file == "" {
		return false
	}

	values, err := readFrontMatter(route.file)

	return err != nil || values.visibility() != res.NoneValue
}

func (server *Server) isActiveFallback(fallback pageRoute) bool {
	if !fallback.isSpecial() || !server.isPublished(fallback) {
		return false
	}

	directory := path.Dir(fmt.Sprintf(res.RootedPathFormat, fallback.source))
	index := server.findPage(path.Join(directory, res.IndexName))
	notFound := server.findPage(path.Join(directory, res.NotFoundName))

	if index.file == "" && notFound.file == "" {
		return true
	}

	_, isWarned := server.warnings.LoadOrStore(fallback.file, true)

	if !isWarned {
		server.logger.Printf(res.IgnoredFallbackMessage, fallback.file)
	}

	return false
}

func (server *Server) buildMissingPage(page missingPage) pageResponse {
	values, body, err := loadPage(page.route.file, page.route.url)

	if err != nil {
		server.logger.Printf(res.FileLogFormat, server.sitePath(page.route.file), err)
		return pageResponse{status: page.status}
	}

	return pageResponse{status: page.status, body: server.render(page.route.file, body, values)}
}
