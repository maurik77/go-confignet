package providers

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/Maurik77/go-confignet/extensions"

	"github.com/lafriks/go-shamir"
)

// SplittedSecretsConfigurationProvider loads configuration from other configuration providers and decrypt the splitted values
type SplittedSecretsConfigurationProvider struct {
	data                   map[string]string
	configurationProviders []extensions.IConfigurationProvider
}

// Add adds the configuration provider to the inner collection
func (provider *SplittedSecretsConfigurationProvider) Add(source extensions.IConfigurationProvider) {
	provider.configurationProviders = append(provider.configurationProviders, source)
	log.Printf("SplittedSecretsConfigurationProvider:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// Load configuration from environment variables
func (provider *SplittedSecretsConfigurationProvider) Load() {
	provider.data = make(map[string]string)
	mapParts := make(map[string][][]byte)

	for _, confProvider := range provider.configurationProviders {
		confProvider.Load()

		for key, value := range confProvider.GetData() {
			decodedString, err := base64.StdEncoding.DecodeString(value)

			if err != nil {
				continue
			}

			internalKey := strings.ReplaceAll(key, confProvider.GetSeparator(), provider.GetSeparator())
			mapParts[internalKey] = append(mapParts[internalKey], decodedString)
		}
	}

	for key, parts := range mapParts {
		decryptedBytes, err := shamir.Combine(parts...)

		if err == nil {
			provider.data[key] = string(decryptedBytes)
		} else {
			log.Printf("SplittedSecretsConfigurationProvider:Unable to decryptkey '%v'. Error: %v", key, err)
		}
	}
}

// GetData provides the loaded data
func (provider *SplittedSecretsConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *SplittedSecretsConfigurationProvider) GetSeparator() string {
	return ":"
}
