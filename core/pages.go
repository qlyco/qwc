package core

import (
	"errors"
	"net/http"

	"github.com/qlyco/qwc/internal/res"
)

func (server *Server) buildPage(route pageRoute) pageResponse {
	values, body, err := server.loadRouteValues(route)

	if err != nil {
		server.logger.Printf(res.FileLogFormat, server.sitePath(route.file), err)
		return pageResponse{status: http.StatusNotFound}
	}

	for _, problem := range values.problems() {
		server.logger.Printf(res.FileLogFormat, server.sitePath(route.file), problem)
	}

	if values.visibility() == res.NoneValue {
		return pageResponse{status: http.StatusNotFound}
	}

	pattern, isPaginated := values.text(res.PaginateKey)

	if route.list != nil {
		isPaginated = true
	}

	if route.pageNumber > 0 && !isPaginated {
		return pageResponse{status: http.StatusNotFound}
	}

	if route.pageNumber == 1 {
		return pageResponse{status: http.StatusMovedPermanently, location: server.address.relativeURL(route.url)}
	}

	if isPaginated {
		pagination, err := server.paginate(pattern, values, route)

		if errors.Is(err, errPageOutOfRange) {
			return pageResponse{status: http.StatusNotFound}
		}

		if err != nil {
			server.logger.Printf(res.FileLogFormat, server.sitePath(route.file), err)
		}

		values[res.PaginationKey] = pagination
	}

	return pageResponse{status: http.StatusOK, body: server.render(route.file, body, values)}
}
