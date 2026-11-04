package core

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/qlyco/qwc/internal/res"
	"golang.org/x/sync/errgroup"
)

func (server *Server) Build(parentDir string, siteURL string) (string, error) {
	address, err := parseSiteAddress(siteURL)

	if err != nil {
		return "", err
	}

	builder := *server
	builder.address = address

	return builder.build(filepath.Join(parentDir, res.SiteDir))
}

func (server *Server) build(siteDir string) (string, error) {
	builtAt := time.Now()
	err := os.RemoveAll(siteDir)

	if err != nil {
		return "", err
	}

	var builds errgroup.Group

	builds.SetLimit(runtime.GOMAXPROCS(0))

	err = server.walkPages(func(route pageRoute) error {
		builds.Go(func() error {
			if route.isSpecial() {
				return server.buildMissingFiles(siteDir, route, builtAt)
			}

			return server.buildRoute(siteDir, route, builtAt)
		})

		return nil
	})

	if err == nil {
		err = walkPaths(server.static, func(relativePath string, filePath string) error {
			builds.Go(func() error {
				return copyFile(filepath.Join(siteDir, filepath.FromSlash(relativePath)), filePath, builtAt)
			})

			return nil
		})
	}

	return siteDir, errors.Join(err, builds.Wait())
}

func (server *Server) buildRoute(siteDir string, route pageRoute, builtAt time.Time) error {
	for pageNumber := 0; ; pageNumber = max(pageNumber+1, 2) {
		route.pageNumber = pageNumber
		response := server.buildPage(route)

		if response.status != http.StatusOK {
			return nil
		}

		err := writeAtomic(filepath.Join(siteDir, filepath.FromSlash(route.sitePath())), bytes.NewReader(response.body), builtAt)

		if err != nil {
			return err
		}
	}
}

func (server *Server) buildMissingFiles(siteDir string, route pageRoute, builtAt time.Time) error {
	if !server.isPublished(route) {
		return nil
	}

	directory := path.Dir(fmt.Sprintf(res.RootedPathFormat, route.source))
	page := missingPage{route: route, status: http.StatusNotFound}
	fileNames := []string{res.NotFoundPage}

	if path.Base(trimExt(route.source)) == res.FallbackName {
		if !server.isActiveFallback(route) {
			return nil
		}

		page.route.url = strings.TrimSuffix(page.route.url, res.FallbackName)
		page.status = http.StatusOK
		fileNames = []string{res.IndexPage, res.NotFoundPage}
	}

	response := server.buildMissingPage(page)

	if response.body == nil {
		return nil
	}

	for _, fileName := range fileNames {
		target := filepath.Join(siteDir, filepath.FromSlash(path.Join(directory, fileName)))
		err := writeAtomic(target, bytes.NewReader(response.body), builtAt)

		if err != nil {
			return err
		}
	}

	return nil
}
