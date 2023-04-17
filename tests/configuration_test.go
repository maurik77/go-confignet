package tests

import (
	"github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestConfigurationProviders(t *testing.T) {
	expected := getJSONExpectedValue()

	myCfg := myConfig{}
	err := confignet.Bind("config", &myCfg)
	assert.Nil(t, err)

	validateObject(t, expected, myCfg)

	subObjConf := subObj{}
	err = confignet.Bind("config/Obj1", &subObjConf)
	assert.Nil(t, err)
	validateSubObject(t, *expected.Obj1, subObjConf)
}

func TestConfigurationProvidersWithEnvVars(t *testing.T) {

	t.Setenv("config__Obj1__PropertyString", "envTest")
	t.Setenv("config__Obj1__PropertyInt64", "2377777")
	t.Setenv("config__Obj1__PropertyInt16", "23")
	t.Setenv("config__Obj1__Time", "2022-01-21T10:00:00Z")
	t.Setenv("config__Obj1__ArrayObj__0__PropertyString", "Modified")
	t.Setenv("config__Obj1__ArrayObj__0__PropertyString", "Modified")
	t.Setenv("config__Obj1__ArrayObj__2__PropertyString", "Created")
	t.Setenv("config__Obj1__ArrayInt__4", "5")
	t.Setenv("config__Obj1__MapObj__99__PropertyString", "Created")
	t.Setenv("config__Obj1__MapObj__99__PropertyInt", "88")
	t.Setenv("config__Obj1__MapObj__99__PropertyBool", "true")

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.AddDefaultConfigurationProviders()
	conf := confBuilder.Build()

	myCfg := myConfig{}
	err := conf.Bind("config", &myCfg)
	assert.Nil(t, err)

	expected := getJSONExpectedValue()

	expected.Obj1.PropertyString = "envTest"
	expected.Obj1.PropertyInt64 = 2377777
	expected.Obj1.PropertyInt16 = 23
	expected.Obj1.ArrayObj[0].PropertyString = "Modified"
	expected.Obj1.ArrayObj = append(expected.Obj1.ArrayObj, subObjItem{PropertyString: "Created"})
	expected.Obj1.MapObj[99] = subObjItem{PropertyString: "Created", PropertyInt: 88, PropertyBool: true}

	timeCfg, _ := time.Parse(time.RFC3339Nano, "2022-01-21T10:00:00Z")
	expected.Obj1.Time = timeCfg

	validateObject(t, expected, myCfg)

	subObjConf := subObj{}
	err = conf.Bind("config/Obj1", &subObjConf)
	assert.Nil(t, err)
	validateSubObject(t, *expected.Obj1, subObjConf)
}
