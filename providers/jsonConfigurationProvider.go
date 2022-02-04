package providers

import (
	"encoding/json"
	"log"

	"github.com/Maurik77/go-confignet/internal"
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
func (provider *JSONConfigurationProvider) Load() {
	var payload map[string]interface{}

	err := internal.UnmarshalFromFile(provider.FilePath, &payload, json.Unmarshal)
	if err != nil {
		log.Println("JSONConfigurationProvider:Error during Unmarshal(): ", err)
	}

	provider.data = internal.LoadProperties(provider.GetSeparator(), payload)
}

// GetData provides the loaded data
func (provider *JSONConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *JSONConfigurationProvider) GetSeparator() string {
	return "."
}
