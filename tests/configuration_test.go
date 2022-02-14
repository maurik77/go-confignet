package tests

import (
	"github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"

	"testing"
	"time"
)

func TestConfigurationProviders(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.AddDefaultConfigurationProviders()
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

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

	timeCfg, _ := time.Parse(time.RFC3339Nano, "2022-01-19T10:00:00Z")
	expected.Obj1.Time = timeCfg

	validateObject(t, expected, myCfg)

	subObjConf := subObj{}
	conf.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected.Obj1, subObjConf)
}

func TestConfigurationProvidersWithEnvVars(t *testing.T) {

	t.Setenv("config__Obj1__PropertyString", "envTest")
	t.Setenv("config__Obj1__PropertyInt64", "2377777")
	t.Setenv("config__Obj1__PropertyInt16", "23")
	t.Setenv("config__Obj1__Time", "2022-01-21T10:00:00Z")

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.AddDefaultConfigurationProviders()
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

	expected := myConfig{
		PropertyInt8: 45,
		Obj1: subObj{
			PropertyString: "envTest",
			PropertyInt:    1,
			PropertyInt8:   2,
			PropertyInt16:  23,
			PropertyInt64:  2377777,
			PropertyBool:   true,
		},
	}

	timeCfg, _ := time.Parse(time.RFC3339Nano, "2022-01-21T10:00:00Z")
	expected.Obj1.Time = timeCfg

	validateObject(t, expected, myCfg)

	subObjConf := subObj{}
	conf.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected.Obj1, subObjConf)
}
