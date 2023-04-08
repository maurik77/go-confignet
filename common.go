package confignet

import (
	"sort"

	"github.com/maurik77/go-confignet/decrypters"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/providers"
)

var (
	configurationSources = make(map[string]extensions.IConfigurationSource)
	decrypterSources     = make(map[string]extensions.IConfigurationDecrypterSource)
)

// RegisterDecrypterSource registers a new decrypter source
func RegisterDecrypterSource(decrypterSource extensions.IConfigurationDecrypterSource) {
	if decrypterSource == nil {
		panic("confignet: Register decrypter source is nil")
	}

	uniqueIdentifier := decrypterSource.GetUniqueIdentifier()

	if _, dup := decrypterSources[uniqueIdentifier]; dup {
		panic("sql: Register called twice for decrypter source with unique identifier " + uniqueIdentifier)
	}

	decrypterSources[uniqueIdentifier] = decrypterSource
}

// RegisterConfigurationSource registers a new configuration source
func RegisterConfigurationSource(configurationSource extensions.IConfigurationSource) {
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
	RegisterConfigurationSource(&providers.CmdLineConfigurationProviderSource{})
	RegisterConfigurationSource(&providers.EnvConfigurationProviderSource{})
	RegisterConfigurationSource(&providers.JSONConfigurationProviderSource{})
	RegisterConfigurationSource(&providers.YamlConfigurationProviderSource{})
	RegisterConfigurationSource(&providers.KeyVaultConfigurationProviderSource{})
	RegisterConfigurationSource(&ChainedConfigurationProviderSource{})
	RegisterDecrypterSource(&decrypters.AesConfigurationDecrypterSource{})
	RegisterDecrypterSource(&decrypters.ShamirConfigurationDecrypterSource{})
}
