package core

import (
	"path"
	"slices"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func (server *Server) collection(pattern string) ([]templateData, error) {
	return server.listPages(pattern, false)
}

func (server *Server) allPages(pattern string) ([]templateData, error) {
	return server.listPages(pattern, true)
}

func (server *Server) listPages(pattern string, withIndexes bool) ([]templateData, error) {
	err := validatePattern(pattern)

	if err != nil {
		return nil, err
	}

	entries := []templateData{}

	err = server.walkPages(func(route pageRoute) error {
		name := trimExt(route.source)
		isIndex := path.Base(name) == res.IndexName

		if route.isText || route.isSpecial() || (isIndex && !withIndexes) || !matchPattern(pattern, name) {
			return nil
		}

		values, err := readFrontMatter(route.file)

		if err != nil {
			return err
		}

		if values.visibility() != res.VisibleValue {
			return nil
		}

		values[res.URLKey] = route.url
		entries = append(entries, values)

		return nil
	})

	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(entries, compareEntries)

	return entries, nil
}

func compareEntries(first templateData, second templateData) int {
	byDate := second.date().Compare(first.date())

	if byDate != 0 {
		return byDate
	}

	firstURL, _ := first.text(res.URLKey)
	secondURL, _ := second.text(res.URLKey)

	return strings.Compare(firstURL, secondURL)
}
