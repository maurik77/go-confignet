package providers

import (
	"log"
	"strings"

	"github.com/maurik77/go-confignet/extensions"
)

// ChainedConfigurationProvider loads configuration from other configuration providers
type ChainedConfigurationProvider struct {
	data                       map[string]string
	dataValues                 map[string][]string
	configurationProvidersInfo []extensions.ConfigurationProviderInfo
}

// Add adds the configuration provider to the inner collection
func (provider *ChainedConfigurationProvider) Add(source extensions.IConfigurationProvider) {
	provider.configurationProvidersInfo = append(provider.configurationProvidersInfo, extensions.ConfigurationProviderInfo{Provider: source})
	log.Printf("ChainedConfigurationProvider:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// AddWithEncrypter adds the configuration provider and the decrypter to the inner collection
func (provider *ChainedConfigurationProvider) AddWithEncrypter(source extensions.IConfigurationProvider, decrypter extensions.IConfigurationDecrypter) {
	provider.configurationProvidersInfo = append(provider.configurationProvidersInfo, extensions.ConfigurationProviderInfo{Provider: source, Decrypter: decrypter})
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n, Decrypter:'%T'", source, source.GetSeparator(), decrypter)
}

// Load configuration from environment variables
func (provider *ChainedConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	provider.data = make(map[string]string)
	provider.dataValues = make(map[string][]string)

	for _, confProvider := range provider.configurationProvidersInfo {
		confProvider.Provider.Load(confProvider.Decrypter)

		for key, value := range confProvider.Provider.GetData() {
			internalKey := strings.ReplaceAll(key, confProvider.Provider.GetSeparator(), provider.GetSeparator())
			provider.data[internalKey] = value
			provider.dataValues[internalKey] = append(provider.dataValues[internalKey], value)
		}
	}

	if decrypter != nil {

		for key, values := range provider.dataValues {
			value, err := decrypter.Decrypt(values...)

			if err != nil {
				log.Printf("ChainedConfigurationProvider:Error calling decryption for key %v. %v", key, err)
			} else {
				provider.data[key] = value
			}
		}
	}
}

// GetData provides the loaded data
func (provider *ChainedConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetDataMultiValues provides the loaded data
func (provider *ChainedConfigurationProvider) GetDataMultiValues() map[string][]string {
	return provider.dataValues
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *ChainedConfigurationProvider) GetSeparator() string {
	return ":"
}
