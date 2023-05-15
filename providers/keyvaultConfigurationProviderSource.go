package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

const (
	// ConfigurationProviderKeyVaultIdentifier is the environment variable containing the UniqueIdentifier of the configuration provider
	ConfigurationProviderKeyVaultIdentifier = "keyvault"
)

// KeyVaultConfigurationProviderSource is able to create KeyVaultConfigurationProvider starting from the provider settings
type KeyVaultConfigurationProviderSource struct {
}

// NewConfigurationProvider creates KeyVaultConfigurationProvider starting from the provider settings
func (providerSource *KeyVaultConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		return nil, fmt.Errorf("KeyVaultConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.GetPropertyValue("prefix", "").(string)
	removePrefix := settings.GetPropertyValue("removePrefix", false).(bool)
	tenantID := settings.GetPropertyValue("tenantID", "").(string)
	clientID := settings.GetPropertyValue("clientID", "").(string)
	clientSecret := settings.GetPropertyValue("clientSecret", "").(string)
	baseURL := settings.GetPropertyValue("baseURL", "").(string)

	return &KeyVaultConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      baseURL,
	}, nil
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *KeyVaultConfigurationProviderSource) GetUniqueIdentifier() string {
	return ConfigurationProviderKeyVaultIdentifier
}
