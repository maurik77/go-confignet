package extensions

// IConfigurationBuilder is the configuration builder interface
type IConfigurationBuilder interface {
	Add(source IConfigurationProvider)
	AddWithEncrypter(source IConfigurationProvider, decrypter IConfigurationDecrypter)
	Build() IConfiguration
	AddDefaultConfigurationProviders()
	AddDefaultConfigurationProvidersWithBasePath(basePath string)
	ConfigureConfigurationProviders()
	ConfigureConfigurationProvidersFromJSONConfig(jsonPath string)
	ConfigureConfigurationProvidersFromYamlConfig(jsonPath string)
	ConfigureConfigurationProvidersFromSettings(settings Settings)
}
