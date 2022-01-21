package providers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// YamlConfigurationProvider loads configuration from YAML file key-value pairs
type YamlConfigurationProvider struct {
	FilePath string
	data     map[string]string
}

// Load configuration from YAML file key-value pairs
func (provider *YamlConfigurationProvider) Load() {
	provider.data = make(map[string]string)

	if provider.FilePath == "" {
		provider.FilePath = "app.yaml"
	}

	if _, err := os.Stat(provider.FilePath); errors.Is(err, os.ErrNotExist) {
		log.Printf("YamlConfigurationProvider:File not found %v", provider.FilePath)
		return
	}

	content, err := ioutil.ReadFile(provider.FilePath)

	if err != nil {
		log.Println("YamlConfigurationProvider:Error when opening file: ", err)
		return
	}

	var payload map[string]interface{}
	err = yaml.Unmarshal(content, &payload)
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
