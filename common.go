package confignet

import (
	"sort"

	"github.com/Maurik77/go-confignet/extensions"
	"github.com/Maurik77/go-confignet/providers"
)

var (
	configurationSources = make(map[string]extensions.IConfigurationSource)
)

//Register registers a new configuration source
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

// ConfigurationSources returns a sorted list of the names of the registered configuration sources.
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
