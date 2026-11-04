package core

import (
	"errors"
	"strconv"

	"github.com/qlyco/qwc/internal/res"
)

var errPageOutOfRange = errors.New(res.PageOutOfRangeMessage)

type paginationInfo struct {
	Items      []templateData
	Page       int
	TotalPages int
	Total      int
	PrevURL    string
	NextURL    string
}

func (server *Server) paginate(pattern string, values templateData, route pageRoute) (*paginationInfo, error) {
	pageSize := server.readPageSize(values)
	entries, err := server.pageEntries(pattern, route)

	if err != nil {
		return nil, err
	}

	pageNumber := max(route.pageNumber, 1)
	totalPages := max(1, (len(entries)+pageSize-1)/pageSize)

	if pageNumber > totalPages {
		return nil, errPageOutOfRange
	}

	start := (pageNumber - 1) * pageSize
	end := min(start+pageSize, len(entries))

	pagination := &paginationInfo{
		Items:      entries[start:end],
		Page:       pageNumber,
		TotalPages: totalPages,
		Total:      len(entries),
	}

	if pageNumber > 1 {
		pagination.PrevURL = pageNumberURL(route.url, pageNumber-1)
	}

	if pageNumber < totalPages {
		pagination.NextURL = pageNumberURL(route.url, pageNumber+1)
	}

	return pagination, nil
}

func (server *Server) readPageSize(values templateData) int {
	size, isSize := parsePageSize(values)

	if isSize {
		return size
	}

	sitewide, err := server.loadSitewideData()

	if err == nil {
		size, isSize = parsePageSize(sitewide)
	}

	if isSize {
		return size
	}

	return res.PageSizeDefault
}

func parsePageSize(values templateData) (int, bool) {
	switch value := values[res.PageSizeKey].(type) {
	case string:
		size, err := strconv.Atoi(value)
		return size, err == nil && size >= 1
	case float64:
		return int(value), value >= 1 && value == float64(int(value))
	}

	return 0, false
}

func (server *Server) pageEntries(pattern string, route pageRoute) ([]templateData, error) {
	if route.list != nil {
		return route.list()
	}

	return server.collection(resolvePattern(pattern, route.source))
}
