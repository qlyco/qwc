package core

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/qlyco/qwc/internal/res"
)

type siteAddress struct {
	origin string
	path   string
}

func parseSiteAddress(address string) (siteAddress, error) {
	if address == "" {
		return siteAddress{}, nil
	}

	if !strings.Contains(address, res.SchemeSeparator) {
		address = fmt.Sprintf(res.DefaultSchemeFormat, address)
	}

	parsed, err := url.Parse(address)

	if err != nil || parsed.Host == "" {
		return siteAddress{}, fmt.Errorf(res.InvalidURLMessage, address)
	}

	return siteAddress{
		origin: fmt.Sprintf(res.OriginFormat, parsed.Scheme, parsed.Host),
		path:   strings.TrimSuffix(parsed.Path, res.PathSeparator),
	}, nil
}

func (address siteAddress) relativeURL(target string) string {
	if !isSitePath(target) {
		return target
	}

	return fmt.Sprintf(res.JoinFormat, address.path, target)
}

func (address siteAddress) absoluteURL(target string) string {
	if !isSitePath(target) {
		return target
	}

	return fmt.Sprintf(res.JoinFormat, address.origin, address.relativeURL(target))
}

func (address siteAddress) sitePath(requestPath string) (string, bool) {
	if address.path == "" || requestPath == address.path {
		return strings.TrimPrefix(requestPath, address.path), true
	}

	return strings.CutPrefix(requestPath, fmt.Sprintf(res.JoinFormat, address.path, res.PathSeparator))
}

func isSitePath(target string) bool {
	return strings.HasPrefix(target, res.PathSeparator) && !strings.HasPrefix(target, res.ProtocolRelativePrefix)
}
