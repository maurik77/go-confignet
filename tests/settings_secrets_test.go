package tests

import (
	"github.com/Maurik77/go-confignet"
	"github.com/Maurik77/go-confignet/extensions"

	"testing"
)

func TestConfigureConfigurationProvidersSecret(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	t.Setenv(confignet.EnvConfigFileType, "json")
	t.Setenv(confignet.EnvConfigFilePath, "settings-secrets.json")

	confBuilder.ConfigureConfigurationProviders()
	config := confBuilder.Build()

	validateBinding(config, t)
}
