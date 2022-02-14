package confignet

import "github.com/maurik77/go-confignet/extensions"

const (
	// ConfigurationProviderChainedIdentifier is the environment variable containing the UniqueIdentifier of the configuration provider
	ConfigurationProviderChainedIdentifier = "chained"
)

// ChainedConfigurationProviderSource is able to create ChainedConfigurationProvider starting from the provider settings
type ChainedConfigurationProviderSource struct {
}

// NewConfigurationProvider creates ChainedConfigurationProvider starting from the provider settings
func (providerSource *ChainedConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("ChainedConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	return &ChainedConfigurationProvider{}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *ChainedConfigurationProviderSource) GetUniqueIdentifier() string {
	return ConfigurationProviderChainedIdentifier
}
