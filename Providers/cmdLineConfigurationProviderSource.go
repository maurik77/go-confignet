package providers

import (
	extensions "confignet/extensions"
)

// CmdLineConfigurationProviderSource
type CmdLineConfigurationProviderSource struct {
}

// Load configuration from commandline arguments
func (providerSource *CmdLineConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("CmdLineConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.Properties["prefix"].(string)
	removePrefix := settings.Properties["removePrefix"].(bool)

	//TODO KeyMapper
	return &CmdLineConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
	}
}

// GetData provides the loaded data
func (providerSource *CmdLineConfigurationProviderSource) GetUniqueIdentifier() string {
	return "cmdline"
}
