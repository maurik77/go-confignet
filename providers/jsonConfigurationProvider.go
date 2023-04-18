package providers

import (
	"encoding/json"

	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
	"github.com/rs/zerolog/log"
)

const (
	// DefaultJSONFile = app.json
	DefaultJSONFile = "app.json"
)

// JSONConfigurationProvider loads configuration from JSON file key-value pairs
type JSONConfigurationProvider struct {
	FilePath string
	data     map[string]string
}

// Load from JSON file key-value pairs
func (provider *JSONConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	var payload map[string]interface{}

	err := internal.UnmarshalFromFile(provider.FilePath, &payload, json.Unmarshal)
	if err != nil {
		log.Err(err).Msg("JSONConfigurationProvider:Error during Unmarshal()")
	}

	provider.data = internal.LoadProperties(provider.GetSeparator(), payload)

	if decrypter != nil {
		var err error

		for key, value := range provider.data {
			value, err = decrypter.Decrypt(value)

			if err != nil {
				log.Err(err).Msgf("JSONConfigurationProvider:Error calling decryption for key %v", key)
			} else {
				provider.data[key] = value
			}
		}
	}
}

// GetData provides the loaded data
func (provider *JSONConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *JSONConfigurationProvider) GetSeparator() string {
	return "."
}
