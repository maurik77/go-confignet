package providers

import (
	extensions "confignet/extensions"
)

// SplittedSecretsConfigurationProviderSource
type SplittedSecretsConfigurationProviderSource struct {
}

// NewConfigurationProvider
func (providerSource *SplittedSecretsConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("SplittedSecretsConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	return &SplittedSecretsConfigurationProvider{}
}

// GetData provides the loaded data
func (providerSource *SplittedSecretsConfigurationProviderSource) GetUniqueIdentifier() string {
	return "shamir"
}
