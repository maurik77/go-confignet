package confignet

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/lafriks/go-shamir"
	// "github.com/lafriks/go-shamir"
)

// SplittedSecretsConfigurationProvider loads configuration from other configuration providers and decrypt the splitted values
type SplittedSecretsConfigurationProvider struct {
	data                   map[string]string
	configurationProviders []IConfigurationProvider
}

// Add adds the configuration provider to the inner collection
func (conf *SplittedSecretsConfigurationProvider) Add(source IConfigurationProvider) {
	conf.configurationProviders = append(conf.configurationProviders, source)
	log.Printf("SecretSplittedConfigurationProvider:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// Load configuration from environment variables
func (conf *SplittedSecretsConfigurationProvider) Load() {
	conf.data = make(map[string]string)
	mapParts := make(map[string][][]byte)

	for _, confProvider := range conf.configurationProviders {
		confProvider.Load()

		for key, value := range confProvider.GetData() {
			decodedString, err := base64.StdEncoding.DecodeString(value)

			if err != nil {
				continue
			}

			internalKey := strings.ReplaceAll(key, confProvider.GetSeparator(), conf.GetSeparator())
			mapParts[internalKey] = append(mapParts[internalKey], decodedString)
		}
	}

	for key, parts := range mapParts {
		decryptedBytes, err := shamir.Combine(parts...)

		if err == nil {
			conf.data[key] = string(decryptedBytes)
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
