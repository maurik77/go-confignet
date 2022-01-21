package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// JSONConfigurationProvider loads configuration from JSON file key-value pairs
type JSONConfigurationProvider struct {
	FilePath string
	data     map[string]string
}

// Load from JSON file key-value pairs
func (provider *JSONConfigurationProvider) Load() {
	provider.data = make(map[string]string)

	if provider.FilePath == "" {
		provider.FilePath = "app.json"
	}

	if _, err := os.Stat(provider.FilePath); errors.Is(err, os.ErrNotExist) {
		log.Printf("JSONConfigurationProvider:File not found %v", provider.FilePath)
		return
	}

	content, err := ioutil.ReadFile(provider.FilePath)

	if err != nil {
		log.Println("JSONConfigurationProvider:Error when opening file: ", err)
		return
	}

	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Println("JSONConfigurationProvider:Error during Unmarshal(): ", err)
	}

	provider.loadProperties("", payload)
}

// GetData provides the loaded data
func (provider *JSONConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *JSONConfigurationProvider) GetSeparator() string {
	return "."
}

func (provider *JSONConfigurationProvider) loadProperties(parent string, json map[string]interface{}) {
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
