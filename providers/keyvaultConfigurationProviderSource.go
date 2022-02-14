package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

// KeyvaultConfigurationProviderSource is able to create KeyvaultConfigurationProvider starting from the provider settings
type KeyvaultConfigurationProviderSource struct {
}

// NewConfigurationProvider creates KeyvaultConfigurationProvider starting from the provider settings
func (providerSource *KeyvaultConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		return nil, fmt.Errorf("KeyvaultConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.GetPropertyValue("prefix", "").(string)
	removePrefix := settings.GetPropertyValue("removePrefix", false).(bool)
	tenantID := settings.GetPropertyValue("tenantID", "").(string)
	clientID := settings.GetPropertyValue("clientID", "").(string)
	clientSecret := settings.GetPropertyValue("clientSecret", "").(string)
	baseURL := settings.GetPropertyValue("baseURL", "").(string)

	return &KeyvaultConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      baseURL,
	}, nil
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *KeyvaultConfigurationProviderSource) GetUniqueIdentifier() string {
	return "keyvault"
}
