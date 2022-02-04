package providers

import (
	"fmt"
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

	provider.loadProperties("", payload)
}

// GetData provides the loaded data
func (provider *YamlConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *YamlConfigurationProvider) GetSeparator() string {
	return "."
}

func (provider *YamlConfigurationProvider) loadProperties(parent string, json map[string]interface{}) {
	for key, value := range json {
		if parent != "" {
			key = fmt.Sprintf("%v%v%v", parent, provider.GetSeparator(), key)
		}

		switch v := value.(type) {
		default:
			provider.data[key] = fmt.Sprint(v)
		case map[string]interface{}:
			provider.loadProperties(key, v)
		}
	}
}
