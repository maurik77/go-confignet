package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

// EnvConfigurationProviderSource is able to create EnvConfigurationProvider starting from the provider settings
type EnvConfigurationProviderSource struct {
}

// NewConfigurationProvider creates EnvConfigurationProvider starting from the provider settings
func (providerSource *EnvConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		return nil, fmt.Errorf("EnvConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.GetPropertyValue("prefix", "").(string)
	removePrefix := settings.GetPropertyValue("removePrefix", false).(bool)

	return &EnvConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
	}, nil
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *EnvConfigurationProviderSource) GetUniqueIdentifier() string {
	return "env"
}
