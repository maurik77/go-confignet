package providers

import (
	"github.com/BurntSushi/toml"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
)

const (
	// DefaultTOMLFile is the default TOML configuration file name
	DefaultTOMLFile = "app.toml"
)

// TomlConfigurationProvider loads configuration from a TOML file.
// It uses "." (dot) as separator for hierarchical configuration.
type TomlConfigurationProvider struct {
	FilePath string // default: "app.toml"
	data     map[string]string
}

// Load reads the TOML file and populates the internal key-value map.
func (provider *TomlConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	provider.data = make(map[string]string)

	var payload map[string]interface{}
	err := internal.UnmarshalFromFile(provider.FilePath, &payload, toml.Unmarshal)

	if err != nil {
		logger.Printf("TomlConfigurationProvider:Error during Unmarshal(): %v", err)
	}

	provider.data = internal.LoadProperties(provider.GetSeparator(), payload)

	if decrypter != nil {
		for key, value := range provider.data {
			decrypted, err := decrypter.Decrypt(value)
			if err != nil {
				logger.Printf("TomlConfigurationProvider:Error calling decryption for key %v. %v", key, err)
			} else {
				provider.data[key] = decrypted
			}
		}
	}
}

// GetData returns the loaded key-value map.
func (provider *TomlConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator returns the separator used for hierarchical keys.
func (provider *TomlConfigurationProvider) GetSeparator() string {
	return "."
}
