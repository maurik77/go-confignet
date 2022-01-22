package confignet

import (
	"os"
	"testing"
	"time"
)

type myConfig struct {
	Obj1         subObj
	PropertyInt8 int8
}

type subObj struct {
	PropertyString string
	PropertyInt    int
	PropertyInt8   int8
	PropertyInt16  int16
	PropertyInt64  int64
	PropertyBool   bool
	Time           time.Time
}

func TestConfigurationProviders(t *testing.T) {
	var confBuilder IConfigurationBuilder = &ConfigurationBuilder{}
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

	os.Setenv("config__Obj1__PropertyString", "envTest")
	os.Setenv("config__Obj1__PropertyInt64", "2377777")
	os.Setenv("config__Obj1__PropertyInt16", "23")
	os.Setenv("config__Obj1__Time", "2022-01-19")

	var confBuilder IConfigurationBuilder = &ConfigurationBuilder{}
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

func TestChainedConfigurationProviders(t *testing.T) {
	var confBuilder IConfigurationBuilder = &ConfigurationBuilder{}
	var chained IConfigurationProvider = &SplittedSecretsConfigurationProvider{}
	confBuilder.Add(chained)
	conf := confBuilder.Build()

	t.Log("First test", conf)
}

func validateSubObject(t *testing.T, expected subObj, result subObj) {
	if result.PropertyBool != expected.PropertyBool {
		t.Log("error should be", expected.PropertyBool, ", but got", result.PropertyBool)
		t.Fail()
	}

	if result.PropertyInt != expected.PropertyInt {
		t.Log("error should be", expected.PropertyInt, ", but got", result.PropertyInt)
		t.Fail()
	}

	if result.PropertyInt8 != expected.PropertyInt8 {
		t.Log("error should be", expected.PropertyInt8, ", but got", result.PropertyInt8)
		t.Fail()
	}

	if result.PropertyInt16 != expected.PropertyInt16 {
		t.Log("error should be", expected.PropertyInt16, ", but got", result.PropertyInt16)
		t.Fail()
	}

	if result.PropertyInt64 != expected.PropertyInt64 {
		t.Log("error should be", expected.PropertyInt64, ", but got", result.PropertyInt64)
		t.Fail()
	}

	if result.Time != expected.Time {
		t.Log("error should be", expected.Time, ", but got", result.Time)
		t.Fail()
	}
}
