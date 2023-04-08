package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

const (
	// ConfigurationProviderYAMLIdentifier is the environment variable containing the UniqueIdentifier of the configuration provider
	ConfigurationProviderYAMLIdentifier = "yaml"
)

// YamlConfigurationProviderSource is able to create YamlConfigurationProvider starting from the provider settings
type YamlConfigurationProviderSource struct {
}

// NewConfigurationProvider creates YamlConfigurationProvider starting from the provider settings
func (providerSource *YamlConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		return nil, fmt.Errorf("YamlConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	filePath := settings.GetPropertyValue("filePath", "").(string)

	return &YamlConfigurationProvider{
		FilePath: filePath,
	}, nil
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *YamlConfigurationProviderSource) GetUniqueIdentifier() string {
	return ConfigurationProviderYAMLIdentifier
}
