package tests

import (
	"testing"

	confignet "github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/providers"
	"github.com/stretchr/testify/assert"
)

func buildTomlConf(filePath string) extensions.IConfiguration {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.Add(&providers.TomlConfigurationProvider{FilePath: filePath})
	return confBuilder.Build()
}

func TestTomlProvider_Bind(t *testing.T) {
	conf := buildTomlConf("app.toml")

	myCfg := myConfig{}
	err := conf.Bind("config", &myCfg)
	assert.Nil(t, err)

	validateObject(t, getJSONExpectedValue(), myCfg)
}

func TestTomlProvider_BindSubSection(t *testing.T) {
	conf := buildTomlConf("app.toml")

	subObjConf := subObj{}
	err := conf.Bind("config/Obj1", &subObjConf)
	assert.Nil(t, err)

	validateSubObject(t, *getJSONExpectedValue().Obj1, subObjConf)
}

func TestTomlProvider_GetValue(t *testing.T) {
	conf := buildTomlConf("app.toml")

	assert.Equal(t, "TestObj1", conf.GetValue("config/Obj1/PropertyString"))
	assert.Equal(t, "45", conf.GetValue("config/PropertyInt8"))
	assert.Equal(t, "true", conf.GetValue("config/Obj1/PropertyBool"))
	assert.Equal(t, "1.5", conf.GetValue("config/Obj1/PropertyFloat32"))
}

func TestTomlProvider_MissingFile(t *testing.T) {
	conf := buildTomlConf("nonexistent.toml")

	cfg := myConfig{}
	err := conf.Bind("config", &cfg)
	// Bind itself does not error on missing file; the provider logs and returns empty data
	assert.Nil(t, err)
	assert.Nil(t, cfg.Obj1)
}

func TestTomlProvider_MetaConfig(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.ConfigureConfigurationProviders(confignet.ConfigFileTypeJSON, "settings-toml.json")
	conf := confBuilder.Build()

	myCfg := myConfig{}
	err := conf.Bind("config", &myCfg)
	assert.Nil(t, err)
	assert.NotNil(t, myCfg.Obj1)
	assert.Equal(t, "TestObj1", myCfg.Obj1.PropertyString)
}
