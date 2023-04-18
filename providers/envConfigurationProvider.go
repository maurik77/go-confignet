package providers

import (
	"os"
	"strings"

	"github.com/maurik77/go-confignet/extensions"
	"github.com/rs/zerolog/log"
)

// EnvConfigurationProvider loads configuration from environment variables
type EnvConfigurationProvider struct {
	Prefix       string
	RemovePrefix bool
	data         map[string]string
}

// Load configuration from environment variables
func (provider *EnvConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
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

		if decrypter != nil {
			var err error
			value, err = decrypter.Decrypt(value)

			if err != nil {
				log.Err(err).Msgf("EnvConfigurationProvider:Error calling decryption for key %v", key)
			}
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
