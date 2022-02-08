package tests

import (
	"github.com/Maurik77/go-confignet"
	"github.com/Maurik77/go-confignet/extensions"

	"testing"
)

func TestConfigureConfigurationProvidersFromEnvVarJson(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	t.Setenv(confignet.EnvConfigFileType, "json")
	t.Setenv(confignet.EnvConfigFilePath, "settings.json")

	confBuilder.ConfigureConfigurationProviders()
	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func TestConfigureConfigurationProvidersFromEnvVarYaml(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	t.Setenv(confignet.EnvConfigFileType, "yaml")
	t.Setenv(confignet.EnvConfigFilePath, "settings.yaml")

	confBuilder.ConfigureConfigurationProviders()

	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func TestConfigureConfigurationProvidersFromJSONConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProvidersFromJSONConfig("")

	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func TestConfigureConfigurationProvidersFromYamlConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProvidersFromYamlConfig("")

	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func validateBinding(config extensions.IConfiguration, t *testing.T, expected *myConfig) {

	myCfg := myConfig{}
	config.Bind("config", &myCfg)

	if expected == nil {
		expected = &myConfig{
			PropertyInt8: 45,
			Obj1: subObj{
				PropertyString: "Encrypted splitted string",
				PropertyInt:    1,
				PropertyInt8:   2,
				PropertyInt16:  3,
				PropertyInt64:  4,
				PropertyBool:   true,
			}}
	}

	validateObject(t, *expected, myCfg)

	subObjConf := subObj{}
	config.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected.Obj1, subObjConf)
}
