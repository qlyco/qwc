package core

import (
	"cmp"
	"net/http"
	"path"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

type pageResponse struct {
	status      int
	contentType string
	body        []byte
	isGzipped   bool
	location    string
	etag        string
}

func (server *Server) errorResponse(filePath string, err error) pageResponse {
	server.logger.Printf(res.FileLogFormat, filePath, err)

	return pageResponse{status: http.StatusInternalServerError, body: []byte(res.RenderErrorMessage)}
}

func textType(name string) (string, bool) {
	contentType, isText := res.TextTypes[strings.ToLower(path.Ext(name))]

	return contentType, isText
}

func writeResponse(writer http.ResponseWriter, request *http.Request, response pageResponse) {
	if response.location != "" {
		http.Redirect(writer, request, response.location, response.status)
		return
	}

	header := writer.Header()
	header.Set(res.ContentTypeHeader, cmp.Or(response.contentType, res.TextTypes[res.PageExt]))
	body := response.body

	if response.isGzipped {
		header.Set(res.VaryHeader, res.AcceptEncodingHeader)
	}

	if response.etag != "" {
		header.Set(res.ETagHeader, response.etag)
		header.Set(res.CacheControlHeader, res.RevalidateCacheControl)
	}

	if response.etag != "" && isNotModified(request, response.etag) {
		writer.WriteHeader(http.StatusNotModified)
		return
	}

	if response.isGzipped {
		if acceptsGzip(request) {
			header.Set(res.ContentEncodingHeader, res.GzipEncoding)
		} else {
			body, _ = decompress(body)
		}
	}

	writer.WriteHeader(response.status)
	writer.Write(body)
}
