package extensions

// IConfiguration is the interface of the configuration
type IConfiguration interface {
	GetProviders() []IConfigurationProvider
	Bind(section string, value interface{})
	GetValue(section string) string
}
