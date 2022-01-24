package confignet

import (
	extensions "confignet/extensions"
	providers "confignet/providers"
	"fmt"
	"log"
)

// ConfigurationBuilder is the concrete implementation
type ConfigurationBuilder struct {
	configurationProviders []extensions.IConfigurationProvider
	decrypters             map[*extensions.IConfigurationProvider]extensions.IConfigurationDecrypter
}

// Add adds the configuration provider to the inner collection
func (conf *ConfigurationBuilder) Add(source extensions.IConfigurationProvider) {
	conf.configurationProviders = append(conf.configurationProviders, source)
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// AddWithEncrypter adds the configuration provider and the decrypter to the inner collection
func (conf *ConfigurationBuilder) AddWithEncrypter(source extensions.IConfigurationProvider, decrypter extensions.IConfigurationDecrypter) {
	conf.configurationProviders = append(conf.configurationProviders, source)
	conf.decrypters[&source] = decrypter
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n, Decrypter:'%T'", source, source.GetSeparator(), decrypter)
}

// AddDefaultConfigurationProviders adds the default configuration providers
func (conf *ConfigurationBuilder) AddDefaultConfigurationProviders() {
	conf.AddDefaultConfigurationProvidersWithBasePath("")
}

// AddDefaultConfigurationProvidersWithBasePath adds the default configuration providers
func (conf *ConfigurationBuilder) AddDefaultConfigurationProvidersWithBasePath(basePath string) {
	conf.Add(&providers.JSONConfigurationProvider{FilePath: fmt.Sprintf("%v%v", basePath, providers.DefaultJSONFile)})
	conf.Add(&providers.YamlConfigurationProvider{FilePath: fmt.Sprintf("%v%v", basePath, providers.DefaultYAMLFile)})
	conf.Add(&providers.EnvConfigurationProvider{})
	conf.Add(&providers.CmdLineConfigurationProvider{})
	conf.Add(&providers.KeyvaultConfigurationProvider{})
}

// Build invokes the load function of each configuration provider and return the Configuration object
func (conf *ConfigurationBuilder) Build() extensions.IConfiguration {
	for _, confProvider := range conf.configurationProviders {
		confProvider.Load()
	}

	result := Configuration{
		configurationProviders: conf.configurationProviders,
		decrypters:             conf.decrypters,
	}

	return &result
}
