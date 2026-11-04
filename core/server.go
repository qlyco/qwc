package core

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/qlyco/qwc/internal/res"
	"golang.org/x/sync/singleflight"
)

var pluginName = regexp.MustCompile(res.PluginNamePattern)

type Server struct {
	root     string
	pages    string
	static   string
	data     string
	sitewide string
	layouts  string
	shared   []string
	cache    cacheDir
	stamp    buildStamp
	logger   *log.Logger
	address  siteAddress
	verbose  bool
	warnings *sync.Map
	plugins  *sync.Map
}

func New(root string, logger *log.Logger) *Server {
	logger = cmp.Or(logger, log.New(io.Discard, "", 0))

	return &Server{
		root:     root,
		pages:    filepath.Join(root, res.PagesDir),
		static:   filepath.Join(root, res.StaticDir),
		data:     filepath.Join(root, res.DataDir),
		sitewide: filepath.Join(root, res.SitewideFile),
		layouts:  filepath.Join(root, res.LayoutsDir),
		shared: []string{
			filepath.Join(root, res.IncludesDir),
			filepath.Join(root, res.ComponentsDir),
		},
		cache:    cacheDir{path: filepath.Join(root, res.CacheDir), logger: logger, builds: &singleflight.Group{}},
		stamp:    buildStamp{path: filepath.Join(root, res.BuildFile), logger: logger},
		logger:   logger,
		warnings: &sync.Map{},
		plugins:  &sync.Map{},
	}
}

func (server *Server) SetVerbose(isVerbose bool) {
	server.verbose = isVerbose
	server.Reload()
}

func (server *Server) RegisterPlugin(name string, plugin any) error {
	if !pluginName.MatchString(name) {
		return fmt.Errorf(res.InvalidPluginNameMessage, name)
	}

	_, isTaken := server.plugins.LoadOrStore(name, plugin)

	if isTaken {
		return fmt.Errorf(res.TakenPluginNameMessage, name)
	}

	return nil
}

func (server *Server) plugin(name string) (any, error) {
	plugin, isFound := server.plugins.Load(name)

	if !isFound {
		return nil, fmt.Errorf(res.UnknownPluginMessage, name)
	}

	return plugin, nil
}

func (server *Server) sitePath(filePath string) string {
	relativePath, err := filepath.Rel(server.root, filePath)

	if err != nil {
		return filepath.ToSlash(filepath.Base(filePath))
	}

	return filepath.ToSlash(relativePath)
}

func (server *Server) SetURL(siteURL string) error {
	address, err := parseSiteAddress(siteURL)

	if err != nil {
		return err
	}

	server.address = address
	server.Reload()

	return nil
}

func (server *Server) Reload() {
	server.stamp.renew()
}

func (server *Server) Refresh() {
	server.cache.clear()
	server.stamp.renew()
}

func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	sitePath, isInside := server.address.sitePath(request.URL.Path)

	if !isInside && request.URL.Path == res.PathSeparator {
		writeResponse(writer, request, pageResponse{status: http.StatusFound, location: server.address.relativeURL(res.PathSeparator)})
		return
	}

	if !isInside {
		writeResponse(writer, request, pageResponse{status: http.StatusNotFound})
		return
	}

	urlPath := path.Clean(fmt.Sprintf(res.RootedPathFormat, sitePath))

	if server.serveStatic(writer, request, urlPath) {
		return
	}

	validFrom := server.stamp.validFrom()
	response := server.servePage(server.findRoute(urlPath), validFrom)

	if response.status == http.StatusNotFound {
		response = server.missingResponse(urlPath, validFrom)
	}

	writeResponse(writer, request, response)
}

func (server *Server) servePage(route pageRoute, validFrom int64) pageResponse {
	if route.file == "" {
		return pageResponse{status: http.StatusNotFound}
	}

	return server.cache.text(route.sitePath(), http.StatusOK, validFrom, func() pageResponse {
		return server.buildPage(route)
	})
}
