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
	ConfigureConfigurationProviders()
	ConfigureConfigurationProvidersFromJSONConfig(jsonPath string)
	ConfigureConfigurationProvidersFromYamlConfig(yamlPath string)
	ConfigureConfigurationProvidersFromSettings(settings Settings)
}
