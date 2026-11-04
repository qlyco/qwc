package res

const (
	ListenAddressFormat = ":%d"
	RootedPathFormat    = "/%s"
	JoinFormat          = "%s%s"
	TemplateNameFormat  = "%s/%s"
	PageNumberURLFormat = "%s/%s/%d"
	FileLogFormat       = "%s: %s"
	FileErrorFormat     = "%s: %w"
	PluginNamePattern   = `^[A-Za-z][A-Za-z0-9]*$`
	OriginFormat        = "%s://%s"
	ETagFormat          = `W/"%d"`
	DefaultSchemeFormat = "https://%s"
	TagURLFormat        = "/%s/%s/"
)
