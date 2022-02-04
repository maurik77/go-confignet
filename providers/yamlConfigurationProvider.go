package providers

import (
	"log"

	"github.com/Maurik77/go-confignet/internal"
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
func (provider *YamlConfigurationProvider) Load() {
	provider.data = make(map[string]string)
	var payload map[string]interface{}
	err := internal.UnmarshalFromFile(provider.FilePath, &payload, yaml.Unmarshal)

	if err != nil {
		log.Println("YamlConfigurationProvider:Error during Unmarshal(): ", err)
	}

	provider.data = internal.LoadProperties(provider.GetSeparator(), payload)
}

// GetData provides the loaded data
func (provider *YamlConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *YamlConfigurationProvider) GetSeparator() string {
	return "."
}
