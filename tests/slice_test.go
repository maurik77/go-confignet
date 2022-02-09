package tests

import (
	"github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/providers"

	"testing"
)

func TestConfigurationSlice(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.Add(&providers.JSONConfigurationProvider{FilePath: "app-slice.json"})
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

	// expected := subObj{}

	// validateSubObject(t, expected, myCfg.Obj1)

	// subObjConf := subObj{}
	// conf.Bind("config/Obj1", &subObjConf)
	// validateSubObject(t, expected, subObjConf)
}
