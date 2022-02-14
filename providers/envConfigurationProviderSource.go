package providers

import "github.com/maurik77/go-confignet/extensions"

const (
	// ConfigurationProviderEnvIdentifier is the environment variable containing the UniqueIdentifier of the configuration provider
	ConfigurationProviderEnvIdentifier = "env"
)

// EnvConfigurationProviderSource is able to create EnvConfigurationProvider starting from the provider settings
type EnvConfigurationProviderSource struct {
}

// NewConfigurationProvider creates EnvConfigurationProvider starting from the provider settings
func (providerSource *EnvConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("EnvConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.GetPropertyValue("prefix", "").(string)
	removePrefix := settings.GetPropertyValue("removePrefix", false).(bool)

	return &EnvConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
	}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *EnvConfigurationProviderSource) GetUniqueIdentifier() string {
	return ConfigurationProviderEnvIdentifier
}
