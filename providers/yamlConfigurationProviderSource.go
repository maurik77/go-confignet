package providers

import "github.com/maurik77/go-confignet/extensions"

// YamlConfigurationProviderSource is able to create YamlConfigurationProvider starting from the provider settings
type YamlConfigurationProviderSource struct {
}

// NewConfigurationProvider creates YamlConfigurationProvider starting from the provider settings
func (providerSource *YamlConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("YamlConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	filePath := settings.GetPropertyValue("filePath", "").(string)

	return &YamlConfigurationProvider{
		FilePath: filePath,
	}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *YamlConfigurationProviderSource) GetUniqueIdentifier() string {
	return "yaml"
}
