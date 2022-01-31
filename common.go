package confignet

import (
	extensions "confignet/extensions"
	providers "confignet/providers"
	"sort"
)

var (
	configurationSources = make(map[string]extensions.IConfigurationSource)
)

func Register(configurationSource extensions.IConfigurationSource) {
	if configurationSource == nil {
		panic("confignet: Register configuration source is nil")
	}

	uniqueIdentifier := configurationSource.GetUniqueIdentifier()

	if _, dup := configurationSources[uniqueIdentifier]; dup {
		panic("sql: Register called twice for configuration source with unique identifier " + uniqueIdentifier)
	}

	configurationSources[uniqueIdentifier] = configurationSource
}

// Drivers returns a sorted list of the names of the registered configuration sources.
func ConfigurationSources() []string {
	list := make([]string, 0, len(configurationSources))
	for name := range configurationSources {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func init() {
	Register(&providers.CmdLineConfigurationProviderSource{})
	Register(&providers.EnvConfigurationProviderSource{})
	Register(&providers.JSONConfigurationProviderSource{})
	Register(&providers.YamlConfigurationProviderSource{})
	Register(&providers.KeyvaultConfigurationProviderSource{})
	Register(&providers.SplittedSecretsConfigurationProviderSource{})
}
