package res

const (
	AppName          = "qwc"
	AppDescription   = "web composer"
	VersionDefault   = "development"
	ModulePath       = "github.com/qlyco/qwc"
	UnstampedVersion = "(devel)"
)

const (
	RootFlag      = "root"
	RootDefault   = "public"
	RootUsage     = "folder that holds the site"
	PortFlag      = "port"
	PortUsage     = "port to serve on"
	BuildFlag     = "build"
	BuildUsage    = "build the site as static files, then exit"
	OutputFlag    = "output"
	OutputUsage   = "folder to build _qwcsite in"
	OutputDefault = "."
	URLFlag       = "url"
	URLUsage      = "address the site is served from, such as example.com or user.github.io/project"
	CertFlag      = "cert"
	CertUsage     = "certificate file, serves HTTPS together with -key"
	KeyFlag       = "key"
	KeyUsage      = "private key file for -cert"
	ConfigFlag    = "config"
	ConfigUsage   = "JSON file with settings, see example.qwcconfig.json"
	VerboseFlag   = "verbose"
	VerboseUsage  = "show full error messages on pages, not only [qwc: error]"
	HTTPScheme    = "http"
	HTTPSScheme   = "https"
)

const (
	PagesDir        = "pages"
	StaticDir       = "static"
	DataDir         = "data"
	LayoutsDir      = "layouts"
	IncludesDir     = "includes"
	ComponentsDir   = "components"
	TagsDir         = "tags"
	SitewideFile    = "site.json"
	BuildFile       = "_qwcbuild"
	CacheDir        = ".qwccache"
	SiteDir         = "_qwcsite"
	TempFilePattern = ".qwc-*.tmp"
	IndexPage       = "index.html"
	IndexName       = "index"
	NotFoundPage    = "404.html"
	NotFoundName    = "404"
	FallbackName    = "200"
	PageExt         = ".html"
	MarkdownExt     = ".md"
	DataExt         = ".json"
	GzipExt         = ".gz"
)

var PageExts = []string{PageExt, MarkdownExt}

var CommandFlags = []string{BuildFlag, ConfigFlag}

var VisibilityValues = []string{VisibleValue, HiddenValue, NoneValue}

var TypeValues = []string{PageValue, TextValue}

const (
	WalkRoot               = "."
	PathSeparator          = "/"
	ProtocolRelativePrefix = "//"
	SchemeSeparator        = "://"
	PageSegment            = "page"
	AnySegments            = "**"
	WordSeparator          = " "
	CurrentDir             = "."
	ParentDir              = ".."
)

const (
	FrontMatterStart  = "<!--"
	FrontMatterEnd    = "-->"
	KeyValueSeparator = ":"
	LineBreak         = "\n"
	RawLineBreak      = "\r\n"
	LeadingSpace      = " \t\r\n"
)

const (
	DateKey        = "date"
	URLKey         = "url"
	PaginateKey    = "paginate"
	PageSizeKey    = "page_size"
	SortKeyKey     = "sort_key"
	TagsSlugKey    = "tags_slug"
	PaginationKey  = "pagination"
	DataKey        = "data"
	LayoutKey      = "layout"
	TemplatesKey   = "templates"
	DisabledValue  = "false"
	ContentBlock   = "content"
	VisibilityKey  = "visibility"
	VisibleValue   = "visible"
	HiddenValue    = "hidden"
	NoneValue      = "none"
	TypeKey        = "type"
	PageValue      = "page"
	TextValue      = "text"
	TagsKey        = "tags"
	TagKey         = "tag"
	TermLayout     = "term"
	TaxonomyLayout = "taxonomy"
)

const (
	SitewideFunc    = "site"
	DictFunc        = "dict"
	SliceFunc       = "slice"
	CollectionFunc  = "collection"
	PluginFunc      = "plugin"
	PagesFunc       = "pages"
	RelativeURLFunc = "relativeURL"
	AbsoluteURLFunc = "absoluteURL"
	TaggedFunc      = "tagged"
	TagsFunc        = "tags"
	UpperCaseFunc   = "upperCase"
	LowerCaseFunc   = "lowerCase"
	TitleCaseFunc   = "titleCase"
	ReplaceFunc     = "replace"
	SplitFunc       = "split"
	JoinFunc        = "join"
)

const (
	ContentTypeHeader     = "Content-Type"
	CacheControlHeader    = "Cache-Control"
	AcceptEncodingHeader  = "Accept-Encoding"
	ContentEncodingHeader = "Content-Encoding"
	VaryHeader            = "Vary"
	ETagHeader            = "ETag"
	IfNoneMatchHeader     = "If-None-Match"
)

const (
	RevalidateCacheControl = "no-cache"
	GzipEncoding           = "gzip"
	RefusedQuality         = "q=0"
	ListSeparator          = ","
	ParameterSeparator     = ";"
	WeakETagPrefix         = "W/"
	AnyETag                = "*"
)

var TextTypes = map[string]string{
	".html":        "text/html; charset=utf-8",
	".htm":         "text/html; charset=utf-8",
	".css":         "text/css; charset=utf-8",
	".js":          "text/javascript; charset=utf-8",
	".mjs":         "text/javascript; charset=utf-8",
	".json":        "application/json",
	".webmanifest": "application/manifest+json",
	".xml":         "application/xml",
	".svg":         "image/svg+xml",
	".txt":         "text/plain; charset=utf-8",
	".md":          "text/markdown; charset=utf-8",
	".csv":         "text/csv; charset=utf-8",
}

const (
	BannerMessage            = "%s: %s [%s]\n"
	CustomBannerMessage      = "%s [%s]\n"
	BuiltOnBannerMessage     = "%s [%s], built on %s %s\n"
	ServingMessage           = "Serving %s on %s://localhost%s\n"
	KeyHintMessage           = "Press r to reload or R to refresh\n"
	UnexpectedArgMessage     = "Unexpected argument %q"
	MissingSiteDirMessage    = "Site folder %q not found"
	CertPairMessage          = "-cert and -key must be given together"
	LoadedCertificateMessage = "Loaded the new certificate from %s"
	KeptCertificateMessage   = "Keeping the current certificate, the new files could not be loaded: %v"
	UnknownSettingMessage    = "unknown setting %q in %s"
	InvalidSettingMessage    = "invalid value %s for setting %q in %s"
	ReloadedMessage          = "Reloaded, cached files rebuild on their next request"
	CacheClearedMessage      = "Cleared %q"
	IgnoredFallbackMessage   = "Ignoring %s because its folder also has an index or 404 page"
	BuiltMessage             = "Built %s\n"
)

const (
	RenderErrorMessage       = "Cannot render page"
	DictPairsMessage         = "dict needs key and value pairs"
	DictKeyMessage           = "dict key %v is not a string"
	SliceTypeMessage         = "slice cannot cut %T, it needs a list or text"
	SliceIndexesMessage      = "slice takes a start and an end at most"
	InvalidDateMessage       = "invalid date %q, so it is ignored. Use YYYY-MM-DD or RFC 3339"
	InvalidVisibilityMessage = "invalid visibility %q, so the page is treated as none. Use visible, hidden or none"
	InvalidTypeMessage       = "invalid type %q, so the page is treated as a page. Use page or text"
	ErrorMarker              = "[qwc: error]"
	ErrorMarkerFormat        = "[qwc: %v]"
	MissingLayoutMessage     = "layout %q not found, so %s is shown without it"
	LayoutLoopMessage        = "layout %q is already in the chain for %s, so the chain stops there"
	InvalidPluginNameMessage = "plugin name %q must be letters and digits, starting with a letter"
	TakenPluginNameMessage   = "plugin name %q is already registered"
	UnknownPluginMessage     = "unknown plugin %q"
	InvalidPageSizeMessage   = "invalid page_size %q, so %d is used. Use a whole number above 0"
	InvalidPatternMessage    = "invalid collection pattern %q"
	InvalidURLMessage        = "invalid site address %q"
	PageOutOfRangeMessage    = "page number out of range"
)
