package tests

import (
	"github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"

	"testing"
)

func TestConfigureConfigurationProvidersSecret(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	t.Setenv(confignet.EnvConfigFileType, "json")
	t.Setenv(confignet.EnvConfigFilePath, "settings-secrets.json")

	confBuilder.ConfigureConfigurationProviders()
	config := confBuilder.Build()

	expected := myConfig{
		PropertyInt8: 45,
		Obj1: subObj{
			PropertyString: "TestObj1",
			PropertyInt:    1,
			PropertyInt8:   2,
			PropertyInt16:  3,
			PropertyInt64:  4,
			PropertyBool:   true,
		},
	}

	validateBinding(config, t, &expected)
}
