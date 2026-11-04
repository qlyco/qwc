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
	perPage := readPerPage(values)
	entries, err := server.collection(pattern)

	if err != nil {
		return nil, err
	}

	pageNumber := max(route.pageNumber, 1)
	totalPages := max(1, (len(entries)+perPage-1)/perPage)

	if pageNumber > totalPages {
		return nil, errPageOutOfRange
	}

	start := (pageNumber - 1) * perPage
	end := min(start+perPage, len(entries))

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

func readPerPage(values templateData) int {
	text, isString := values.text(res.PerPageKey)

	if !isString {
		return res.PerPageDefault
	}

	perPage, err := strconv.Atoi(text)

	if err != nil || perPage < 1 {
		return res.PerPageDefault
	}

	return perPage
}
