package confignet

import (
	"log"
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
	for _, confProvider := range conf.configurationProviders {
		confProvider.Load()
	}

	// TODO Retrieve all values with the same key (parts)
	// use shamir to calculate the decrypted value
	// value, err = shamir.Combine(parts...)
}

// GetData provides the loaded data
func (provider *SplittedSecretsConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *SplittedSecretsConfigurationProvider) GetSeparator() string {
	return ":"
}
