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

	expected := subObj{
		PropertyString: "TestObj1",
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
	conf.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected, subObjConf)
}

func TestConfigurationProvidersWithEnvVars(t *testing.T) {

	t.Setenv("config__Obj1__PropertyString", "envTest")
	t.Setenv("config__Obj1__PropertyInt64", "2377777")
	t.Setenv("config__Obj1__PropertyInt16", "23")
	t.Setenv("config__Obj1__Time", "2022-01-19")

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.AddDefaultConfigurationProviders()
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

	expected := subObj{
		PropertyString: "envTest",
		PropertyInt:    1,
		PropertyInt8:   2,
		PropertyInt16:  23,
		PropertyInt64:  2377777,
		PropertyBool:   true,
	}
	timeCfg, _ := time.Parse(time.RFC3339Nano, "2022-01-19")
	expected.Time = timeCfg

	if myCfg.PropertyInt8 != 45 {
		t.Log("error should be", 45, ", but got", myCfg.PropertyInt8)
		t.Fail()
	}

	validateSubObject(t, expected, myCfg.Obj1)

	subObjConf := subObj{}
	conf.Bind("config/Obj1", &subObjConf)
	validateSubObject(t, expected, subObjConf)
}
