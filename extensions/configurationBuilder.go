package extensions

// IConfigurationProviderCollection is the configuration provider collection interface
type IConfigurationProviderCollection interface {
	Add(source IConfigurationProvider)
	AddWithEncrypter(source IConfigurationProvider, decrypter IConfigurationDecrypter)
}

// IConfigurationBuilder is the configuration builder interface
type IConfigurationBuilder interface {
	IConfigurationProviderCollection
	Build() IConfiguration
	AddDefaultConfigurationProviders()
	AddDefaultConfigurationProvidersWithBasePath(basePath string)
	ConfigureConfigurationProvidersFromEnv()
	ConfigureConfigurationProviders(configType string, jsonPath string)
	ConfigureConfigurationProvidersFromSettings(settings Settings)
}
