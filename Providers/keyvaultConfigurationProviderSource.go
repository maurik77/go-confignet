package providers

import (
	extensions "confignet/extensions"
)

// KeyvaultConfigurationProviderSource
type KeyvaultConfigurationProviderSource struct {
}

// Load configuration from commandline arguments
func (providerSource *KeyvaultConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("CmdLineConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.Properties["prefix"].(string)
	removePrefix := settings.Properties["removePrefix"].(bool)
	tenantID := settings.Properties["tenantID"].(string)
	clientID := settings.Properties["clientID"].(string)
	clientSecret := settings.Properties["clientSecret"].(string)
	baseURL := settings.Properties["baseURL"].(string)

	return &KeyvaultConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
		TenantID:     tenantID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      baseURL,
	}
}

// GetData provides the loaded data
func (providerSource *KeyvaultConfigurationProviderSource) GetUniqueIdentifier() string {
	return "keyvault"
}
