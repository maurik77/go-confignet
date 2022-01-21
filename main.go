package main

import (
	"confignet/confignet"
	"fmt"
	"time"
)

type myConfig struct {
	Obj1 subObj
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

func main() {
	var confBuilder confignet.IConfigurationBuilder = &confignet.ConfigurationBuilder{}
	confBuilder.AddDefaultConfigurationProviders()
	// confBuilder.Add(&providers.JsonConfigurationProvider{})
	// confBuilder.Add(&providers.YamlConfigurationProvider{})
	// confBuilder.Add(&providers.EnvConfigurationProvider{Prefix: "config", RemovePrefix: false})
	// confBuilder.Add(&providers.KeyvaultConfigurationProvider{Prefix: "config", RemovePrefix: false})

	conf := confBuilder.Build()

	// for _, p := range conf.GetProviders() {
	// 	fmt.Printf("%T, Separator:'%v'\n", p, p.GetSeparator())
	// 	fmt.Println("Data:{")
	// 	for key, value := range p.GetData() {
	// 		fmt.Printf("\t%v:'%v'\n", key, value)
	// 	}
	// 	fmt.Println("}")
	// 	fmt.Println()
	// }

	myCfg := myConfig{}
	conf.Bind("config", &myCfg)
	fmt.Printf("%+v\n", myCfg)
	fmt.Printf("PropertyString:%v\n", myCfg.Obj1.PropertyString)

	subObj := subObj{}
	conf.Bind("config/Obj1", &subObj)
	fmt.Printf("%+v\n", subObj)
	fmt.Printf("PropertyString:%v\n", subObj.PropertyString)
}
