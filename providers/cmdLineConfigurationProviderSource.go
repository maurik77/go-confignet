package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

const (
	// ConfigurationProviderCmdlineIdentifier is the environment variable containing the UniqueIdentifier of the configuration provider
	ConfigurationProviderCmdlineIdentifier = "cmdline"
)

// CmdLineConfigurationProviderSource is able to create CmdLineConfigurationProvider starting from the provider settings
type CmdLineConfigurationProviderSource struct {
}

// NewConfigurationProvider creates CmdLineConfigurationProvider starting from the provider settings
func (providerSource *CmdLineConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != providerSource.GetUniqueIdentifier() {
		return nil, fmt.Errorf("CmdLineConfigurationProviderSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + providerSource.GetUniqueIdentifier())
	}

	prefix := settings.GetPropertyValue("prefix", "").(string)
	removePrefix := settings.GetPropertyValue("removePrefix", false).(bool)

	// TODO KeyMapper
	return &CmdLineConfigurationProvider{
		Prefix:       prefix,
		RemovePrefix: removePrefix,
	}, nil
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (providerSource *CmdLineConfigurationProviderSource) GetUniqueIdentifier() string {
	return "cmdline"
}
