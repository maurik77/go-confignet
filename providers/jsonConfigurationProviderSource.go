package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

// JSONConfigurationProviderSource is able to create JSONConfigurationProvider starting from the provider settings
type JSONConfigurationProviderSource struct {
}

// NewConfigurationProvider creates JSONConfigurationProvider starting from the provider settings
func (providerSource *JSONConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		return nil, fmt.Errorf("JSONConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	filePath := settings.GetPropertyValue("filePath", "").(string)

	return &JSONConfigurationProvider{
		FilePath: filePath,
	}, nil
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *JSONConfigurationProviderSource) GetUniqueIdentifier() string {
	return "json"
}
