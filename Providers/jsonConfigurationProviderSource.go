package providers

import (
	extensions "confignet/extensions"
)

// JSONConfigurationProviderSource
type JSONConfigurationProviderSource struct {
}

// Load configuration from commandline arguments
func (providerSource *JSONConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("JSONConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	filePath := settings.GetPropertyValue("filePath", "").(string)

	return &JSONConfigurationProvider{
		FilePath: filePath,
	}
}

// GetData provides the loaded data
func (providerSource *JSONConfigurationProviderSource) GetUniqueIdentifier() string {
	return "json"
}
