package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

func etagFor(validFrom int64) string {
	return fmt.Sprintf(res.ETagFormat, validFrom)
}

func isNotModified(request *http.Request, etag string) bool {
	for _, candidate := range strings.Split(request.Header.Get(res.IfNoneMatchHeader), res.ListSeparator) {
		candidate = strings.TrimSpace(candidate)

		if candidate == res.AnyETag || strings.TrimPrefix(candidate, res.WeakETagPrefix) == strings.TrimPrefix(etag, res.WeakETagPrefix) {
			return true
		}
	}

	return false
}
