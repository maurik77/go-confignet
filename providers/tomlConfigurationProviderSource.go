package providers

import (
	"fmt"

	"github.com/maurik77/go-confignet/extensions"
)

const (
	// ConfigurationProviderTOMLIdentifier is the unique identifier for the TOML provider in settings files
	ConfigurationProviderTOMLIdentifier = "toml"
)

// TomlConfigurationProviderSource creates TomlConfigurationProvider instances from provider settings.
type TomlConfigurationProviderSource struct{}

// NewConfigurationProvider creates a TomlConfigurationProvider from the given settings.
func (s *TomlConfigurationProviderSource) NewConfigurationProvider(settings extensions.ProviderSettings) (extensions.IConfigurationProvider, error) {
	if settings.Name != s.GetUniqueIdentifier() {
		return nil, fmt.Errorf("TomlConfigurationProviderSource: settings of configuration source %s has been passed to the configuration source with unique identifier %s", settings.Name, s.GetUniqueIdentifier())
	}

	filePath := settings.GetPropertyValue("filePath", "").(string)

	return &TomlConfigurationProvider{
		FilePath: filePath,
	}, nil
}

// GetUniqueIdentifier returns the identifier used to reference this provider in settings files.
func (s *TomlConfigurationProviderSource) GetUniqueIdentifier() string {
	return ConfigurationProviderTOMLIdentifier
}
