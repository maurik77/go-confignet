package providers

import (
	extensions "confignet/extensions"
)

// SplittedSecretsConfigurationProviderSource is able to create SplittedSecretsConfigurationProvider starting from the provider settings
type SplittedSecretsConfigurationProviderSource struct {
}

// NewConfigurationProvider creates SplittedSecretsConfigurationProvider starting from the provider settings
func (providerSource *SplittedSecretsConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("SplittedSecretsConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	return &SplittedSecretsConfigurationProvider{}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *SplittedSecretsConfigurationProviderSource) GetUniqueIdentifier() string {
	return "shamir"
}
