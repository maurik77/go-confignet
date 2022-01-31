package providers

import (
	extensions "confignet/extensions"
)

// YamlConfigurationProviderSource
type YamlConfigurationProviderSource struct {
}

// Load configuration from commandline arguments
func (providerSource *YamlConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("CmdLineConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	filePath := settings.GetPropertyValue("filePath", "").(string)

	return &YamlConfigurationProvider{
		FilePath: filePath,
	}
}

// GetData provides the loaded data
func (providerSource *YamlConfigurationProviderSource) GetUniqueIdentifier() string {
	return "yaml"
}
