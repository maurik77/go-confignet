package confignet

import (
	providers "confignet/Providers"
	"log"
)

// IConfigurationBuilder is the configuration builder interface
type IConfigurationBuilder interface {
	Add(source IConfigurationProvider)
	Build() IConfiguration
	AddDefaultConfigurationProviders()
}

// ConfigurationBuilder is the concrete implementation
type ConfigurationBuilder struct {
	configurationProviders []IConfigurationProvider
}

// Add adds the configuration provider to the inner collection
func (conf *ConfigurationBuilder) Add(source IConfigurationProvider) {
	conf.configurationProviders = append(conf.configurationProviders, source)
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// AddDefaultConfigurationProviders adds the default configuration providers
func (conf *ConfigurationBuilder) AddDefaultConfigurationProviders() {
	conf.Add(&providers.JSONConfigurationProvider{})
	conf.Add(&providers.YamlConfigurationProvider{})
	conf.Add(&providers.EnvConfigurationProvider{})
	conf.Add(&providers.CmdLineConfigurationProvider{})
	conf.Add(&providers.KeyvaultConfigurationProvider{})
}

// Build invokes the load function of each configuration provider and return the Configuration object
func (conf *ConfigurationBuilder) Build() IConfiguration {
	for _, confProvider := range conf.configurationProviders {
		confProvider.Load()
	}

	result := Configuration{
		configurationProviders: conf.configurationProviders,
	}

	return &result
}
