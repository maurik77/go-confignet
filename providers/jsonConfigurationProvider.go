package providers

import (
	"encoding/json"
	"log"

	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
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
		log.Println("JSONConfigurationProvider:Error during Unmarshal(): ", err)
	}

	provider.data = internal.LoadProperties(provider.GetSeparator(), payload)

	if decrypter != nil {
		var err error

		for key, value := range provider.data {
			value, err = decrypter.Decrypt(value)

			if err != nil {
				log.Printf("JSONConfigurationProvider:Error calling decryption for key %v. %v", key, err)
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
