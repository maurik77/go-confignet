package providers

import (
	"os"
	"strings"
)

// EnvConfigurationProvider loads configuration from environment variables
type EnvConfigurationProvider struct {
	Prefix       string
	RemovePrefix bool
	data         map[string]string
}

// Load configuration from environment variables
func (provider *EnvConfigurationProvider) Load() {
	provider.data = make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.Split(env, "=")
		key := parts[0]
		value := parts[1]

		if provider.Prefix != "" && !strings.HasPrefix(key, provider.Prefix) {
			continue
		}

		if provider.Prefix != "" && provider.RemovePrefix {
			key = strings.TrimPrefix(key, provider.Prefix)
		}

		provider.data[key] = value
	}
}

// GetData provides the loaded data
func (provider *EnvConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *EnvConfigurationProvider) GetSeparator() string {
	return "__"
}
