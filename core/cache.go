package core

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/qlyco/qwc/internal/res"
	"golang.org/x/sync/singleflight"
)

type cacheDir struct {
	path   string
	logger *log.Logger
	builds *singleflight.Group
}

func (cache cacheDir) text(key string, status int, validFrom int64, build func() pageResponse) pageResponse {
	cacheFile := cache.filePath(fmt.Sprintf(res.JoinFormat, key, res.GzipExt))
	contentType, _ := textType(key)

	if isFresh(cacheFile, validFrom) {
		body, err := os.ReadFile(cacheFile)

		if err == nil {
			return pageResponse{status: status, contentType: contentType, body: body, isGzipped: true, etag: cachedETag(status, validFrom)}
		}
	}

	result, _, _ := cache.builds.Do(cacheFile, func() (any, error) {
		return cache.build(cacheFile, status, contentType, validFrom, build), nil
	})

	return result.(pageResponse)
}

func (cache cacheDir) build(cacheFile string, status int, contentType string, validFrom int64, build func() pageResponse) pageResponse {
	startedAt := time.Now()
	response := build()

	if response.status != status || response.body == nil {
		os.Remove(cacheFile)
		return response
	}

	body, err := compress(response.body)

	if err != nil {
		cache.logger.Printf(res.FileLogFormat, cacheFile, err)
		return response
	}

	err = writeAtomic(cacheFile, bytes.NewReader(body), startedAt)

	if err != nil {
		cache.logger.Printf(res.FileLogFormat, cacheFile, err)
	}

	return pageResponse{status: status, contentType: contentType, body: body, isGzipped: true, etag: cachedETag(status, validFrom)}
}

func (cache cacheDir) blob(key string, sourceFile string, validFrom int64) (string, error) {
	cacheFile := cache.filePath(key)

	if isFresh(cacheFile, validFrom) {
		return cacheFile, nil
	}

	_, err, _ := cache.builds.Do(cacheFile, func() (any, error) {
		return nil, copyFile(cacheFile, sourceFile, time.Now())
	})

	return cacheFile, err
}

func (cache cacheDir) filePath(key string) string {
	return filepath.Join(cache.path, filepath.FromSlash(key))
}

func (cache cacheDir) clear() {
	err := os.RemoveAll(cache.path)

	if err != nil {
		cache.logger.Printf(res.FileLogFormat, cache.path, err)
		return
	}

	cache.logger.Printf(res.CacheClearedMessage, cache.path)
}

func isFresh(cacheFile string, validFrom int64) bool {
	info, err := os.Stat(cacheFile)

	return err == nil && info.ModTime().UnixNano() >= validFrom
}

func cachedETag(status int, validFrom int64) string {
	if status != http.StatusOK {
		return ""
	}

	return etagFor(validFrom)
}
