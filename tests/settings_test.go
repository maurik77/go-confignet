package tests

import (
	"time"

	"github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"

	"testing"
)

func TestConfigureConfigurationProvidersFromEnvVarJson(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	t.Setenv(confignet.EnvConfigFileType, confignet.ConfigFileTypeJSON)
	t.Setenv(confignet.EnvConfigFilePath, confignet.DefaultConfigFileJSON)

	confBuilder.ConfigureConfigurationProvidersFromEnv()
	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func TestConfigureConfigurationProvidersFromEnvVarYaml(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	t.Setenv(confignet.EnvConfigFileType, confignet.ConfigFileTypeYAML)
	t.Setenv(confignet.EnvConfigFilePath, confignet.DefaultConfigFileYAML)

	confBuilder.ConfigureConfigurationProvidersFromEnv()

	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func TestConfigureConfigurationProvidersFromJSONConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProviders(confignet.ConfigFileTypeJSON, "")

	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func TestConfigureConfigurationProvidersFromYamlConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProviders(confignet.ConfigFileTypeYAML, "")

	config := confBuilder.Build()

	validateBinding(config, t, nil)
}

func validateBinding(config extensions.IConfiguration, t *testing.T, expected *myConfig) {

	myCfg := myConfig{}
	config.Bind("config", &myCfg)

	if expected == nil {
		jsonFulConfig := getJSONExpectedValue()
		expected = &jsonFulConfig
		expected.PropertyInt8 = 45
		expected.Obj1.PropertyString = "Encrypted splitted string"
		expected.Obj1.PropertyInt = 1
		expected.Obj1.PropertyInt8 = 2
		expected.Obj1.PropertyInt16 = 3
		expected.Obj1.PropertyInt64 = 4
		expected.Obj1.PropertyBool = true

		timeCfg, _ := time.Parse(time.RFC3339Nano, "2022-01-19T10:00:00Z")
		expected.Obj1.Time = timeCfg
	}

	validateObject(t, *expected, myCfg)

	subObjConf := subObj{}
	config.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected.Obj1, subObjConf)
}
