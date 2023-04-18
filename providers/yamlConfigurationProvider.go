package providers

import (
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultYAMLFile = app.json
	DefaultYAMLFile = "app.yaml"
)

// YamlConfigurationProvider loads configuration from YAML file key-value pairs
type YamlConfigurationProvider struct {
	FilePath string
	data     map[string]string
}

// Load configuration from YAML file key-value pairs
func (provider *YamlConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	provider.data = make(map[string]string)
	var payload map[string]interface{}
	err := internal.UnmarshalFromFile(provider.FilePath, &payload, yaml.Unmarshal)

	if err != nil {
		log.Err(err).Msg("YamlConfigurationProvider:Error during Unmarshal()")
	}

	provider.data = internal.LoadProperties(provider.GetSeparator(), payload)

	if decrypter != nil {
		var err error

		for key, value := range provider.data {
			value, err = decrypter.Decrypt(value)

			if err != nil {
				log.Err(err).Msgf("YamlConfigurationProvider:Error calling decryption for key %v", key)
			} else {
				provider.data[key] = value
			}
		}
	}
}

// GetData provides the loaded data
func (provider *YamlConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *YamlConfigurationProvider) GetSeparator() string {
	return "."
}
