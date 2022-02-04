package extensions

// IConfiguration is the interface of the configuration
type IConfiguration interface {
	GetProviders() []ConfigurationProviderInfo
	Bind(section string, value interface{})
	GetValue(section string) string
}

// ConfigurationProviderInfo contains the Configuration Provider and the Decrypter
type ConfigurationProviderInfo struct {
	Provider  IConfigurationProvider
	Decrypter IConfigurationDecrypter
}
