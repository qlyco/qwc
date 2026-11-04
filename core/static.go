package core

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func (server *Server) serveStatic(writer http.ResponseWriter, request *http.Request, urlPath string) bool {
	sourceFile := filepath.Join(server.static, filepath.FromSlash(urlPath))

	if !isFile(sourceFile) {
		return false
	}

	key := strings.TrimPrefix(urlPath, res.PathSeparator)
	contentType, isText := textType(key)
	validFrom := server.stamp.validFrom()

	if isText {
		writeResponse(writer, request, server.cache.text(key, http.StatusOK, validFrom, func() pageResponse {
			return server.readStatic(sourceFile, contentType)
		}))

		return true
	}

	cacheFile, err := server.cache.blob(key, sourceFile, validFrom)

	if err != nil {
		server.logger.Printf(res.FileLogFormat, cacheFile, err)
		cacheFile = sourceFile
	}

	writer.Header().Set(res.ETagHeader, etagFor(validFrom))
	writer.Header().Set(res.CacheControlHeader, res.RevalidateCacheControl)
	http.ServeFile(writer, request, cacheFile)

	return true
}

func (server *Server) readStatic(sourceFile string, contentType string) pageResponse {
	body, err := os.ReadFile(sourceFile)

	if err != nil {
		return server.errorResponse(sourceFile, err)
	}

	return pageResponse{status: http.StatusOK, contentType: contentType, body: body}
}
