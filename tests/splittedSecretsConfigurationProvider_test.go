package tests

import (
	"confignet"
	extensions "confignet/extensions"
	providers "confignet/providers"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/lafriks/go-shamir"
)

func TestGenerateStringParts(t *testing.T) {

	stringByteArray := []byte("Encrytped splitted string")
	cryptParts, _ := shamir.Split(stringByteArray, 3, 2)

	for index, crypt := range cryptParts {
		str := base64.StdEncoding.EncodeToString(crypt)
		fmt.Println(index, str)
		t.Log(index, str)
	}
}

func TestConfigShamir12(t *testing.T) {

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	var chained extensions.IChainedConfigurationProvider = &providers.SplittedSecretsConfigurationProvider{}
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "copy-shamir-1.yaml"})
	chained.Add(&providers.JSONConfigurationProvider{FilePath: "copy-shamir-2.json"})
	confBuilder.Add(chained)
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

	expected := subObj{
		PropertyString: "Encrytped splitted string",
	}

	validateSubObject(t, expected, myCfg.Obj1)
}

func TestConfigShamir13(t *testing.T) {

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	var chained extensions.IChainedConfigurationProvider = &providers.SplittedSecretsConfigurationProvider{}
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "copy-shamir-1.yaml"})
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "copy-shamir-3.yaml"})
	confBuilder.Add(chained)
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

	expected := subObj{
		PropertyString: "Encrytped splitted string",
	}

	validateSubObject(t, expected, myCfg.Obj1)
}

func TestConfigShamir23(t *testing.T) {

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	var chained extensions.IChainedConfigurationProvider = &providers.SplittedSecretsConfigurationProvider{}
	chained.Add(&providers.JSONConfigurationProvider{FilePath: "copy-shamir-2.json"})
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "copy-shamir-3.yaml"})
	confBuilder.Add(chained)
	conf := confBuilder.Build()

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)

	expected := subObj{
		PropertyString: "Encrytped splitted string",
	}

	validateSubObject(t, expected, myCfg.Obj1)
}
