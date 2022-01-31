package tests

import (
	"confignet"
	extensions "confignet/extensions"
	"os"

	"testing"
)

func TestConfigureConfigurationProvidersFromEnvVarJson(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	os.Setenv(confignet.EnvConfigFileType, "json")
	os.Setenv(confignet.EnvConfigFilePath, "settings.json")

	confBuilder.ConfigureConfigurationProviders()
	config := confBuilder.Build()

	validateBinding(config, t)
}

func TestConfigureConfigurationProvidersFromEnvVarYaml(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	os.Setenv(confignet.EnvConfigFileType, "yaml")
	os.Setenv(confignet.EnvConfigFilePath, "settings.yaml")

	confBuilder.ConfigureConfigurationProviders()

	config := confBuilder.Build()

	validateBinding(config, t)
}

func TestConfigureConfigurationProvidersFromJSONConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProvidersFromJSONConfig("")

	config := confBuilder.Build()

	validateBinding(config, t)
}

func TestConfigureConfigurationProvidersFromYamlConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProvidersFromYamlConfig("")

	config := confBuilder.Build()

	validateBinding(config, t)
}

func validateBinding(config extensions.IConfiguration, t *testing.T) {

	myCfg := myConfig{}
	config.Bind("config", &myCfg)

	expected := subObj{
		PropertyString: "Encrypted splitted string",
		PropertyInt:    1,
		PropertyInt8:   2,
		PropertyInt16:  3,
		PropertyInt64:  4,
		PropertyBool:   true,
	}

	if myCfg.PropertyInt8 != 45 {
		t.Log("error should be", 45, ", but got", myCfg.PropertyInt8)
		t.Fail()
	}

	validateSubObject(t, expected, myCfg.Obj1)

	subObjConf := subObj{}
	config.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected, subObjConf)
}
