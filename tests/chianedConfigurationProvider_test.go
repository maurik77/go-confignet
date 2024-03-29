package tests

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/maurik77/go-confignet"
	"github.com/maurik77/go-confignet/decrypters"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/providers"
	"github.com/stretchr/testify/assert"

	"github.com/lafriks/go-shamir"
)

func TestGenerateStringParts(t *testing.T) {

	stringByteArray := []byte("Encrypted splitted string")
	cryptParts, _ := shamir.Split(stringByteArray, 3, 2)

	for index, crypt := range cryptParts {
		str := base64.StdEncoding.EncodeToString(crypt)
		fmt.Println(index, str)
		t.Log(index, str)
	}
}

func TestConfigShamir12(t *testing.T) {

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	var chained extensions.IChainedConfigurationProvider = &confignet.ChainedConfigurationProvider{}
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "shamir/copy-shamir-1.yaml"})
	chained.Add(&providers.JSONConfigurationProvider{FilePath: "shamir/copy-shamir-2.json"})
	confBuilder.AddWithEncrypter(chained, &decrypters.ShamirConfigurationDecrypter{})
	conf := confBuilder.Build()

	myCfg := myConfig{}
	err := conf.Bind("config", &myCfg)

	assert.Nil(t, err)

	expected := subObj{
		PropertyString: "Encrypted splitted string",
	}

	validateSubObject(t, expected, *myCfg.Obj1)
}

func TestConfigShamir13(t *testing.T) {

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	var chained extensions.IChainedConfigurationProvider = &confignet.ChainedConfigurationProvider{}
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "shamir/copy-shamir-1.yaml"})
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "shamir/copy-shamir-3.yaml"})
	confBuilder.AddWithEncrypter(chained, &decrypters.ShamirConfigurationDecrypter{})
	conf := confBuilder.Build()

	myCfg := myConfig{}
	err := conf.Bind("config", &myCfg)
	assert.Nil(t, err)

	expected := subObj{
		PropertyString: "Encrypted splitted string",
	}

	validateSubObject(t, expected, *myCfg.Obj1)
}

func TestConfigShamir23(t *testing.T) {

	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	var chained extensions.IChainedConfigurationProvider = &confignet.ChainedConfigurationProvider{}
	chained.Add(&providers.JSONConfigurationProvider{FilePath: "shamir/copy-shamir-2.json"})
	chained.Add(&providers.YamlConfigurationProvider{FilePath: "shamir/copy-shamir-3.yaml"})
	confBuilder.AddWithEncrypter(chained, &decrypters.ShamirConfigurationDecrypter{})
	conf := confBuilder.Build()

	myCfg := myConfig{}
	err := conf.Bind("config", &myCfg)
	assert.Nil(t, err)

	expected := subObj{
		PropertyString: "Encrypted splitted string",
	}

	validateSubObject(t, expected, *myCfg.Obj1)
}
