package providers

import (
	"log"
	"strings"

	"github.com/Maurik77/go-confignet/extensions"
)

// ChainedConfigurationProvider loads configuration from other configuration providers
type ChainedConfigurationProvider struct {
	data                   map[string]string
	dataValues             map[string][]string
	configurationProviders []extensions.IConfigurationProvider
}

// Add adds the configuration provider to the inner collection
func (provider *ChainedConfigurationProvider) Add(source extensions.IConfigurationProvider) {
	provider.configurationProviders = append(provider.configurationProviders, source)
	log.Printf("ChainedConfigurationProvider:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// Load configuration from environment variables
func (provider *ChainedConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	provider.data = make(map[string]string)
	provider.dataValues = make(map[string][]string)

	for _, confProvider := range provider.configurationProviders {
		confProvider.Load(nil)

		for key, value := range confProvider.GetData() {
			internalKey := strings.ReplaceAll(key, confProvider.GetSeparator(), provider.GetSeparator())
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
