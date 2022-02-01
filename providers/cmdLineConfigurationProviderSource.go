package providers

import "github.com/Maurik77/go-confignet/extensions"

// CmdLineConfigurationProviderSource is able to create CmdLineConfigurationProvider starting from the provider settings
type CmdLineConfigurationProviderSource struct {
}

// NewConfigurationProvider creates CmdLineConfigurationProvider starting from the provider settings
func (providerSource *CmdLineConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) extensions.IConfigurationProvider {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		panic("CmdLineConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.GetPropertyValue("prefix", "").(string)
	removePrefix := settings.GetPropertyValue("removePrefix", false).(bool)

	//TODO KeyMapper
	return &CmdLineConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
	}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *CmdLineConfigurationProviderSource) GetUniqueIdentifier() string {
	return "cmdline"
}
