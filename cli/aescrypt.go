package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/maurik77/go-confignet/internal"
	"gopkg.in/yaml.v2"
)

func main() {
	var configFile, configFileType, secret, destination string
	flag.StringVar(&configFile, "config", "config.yaml", "settings file path")
	flag.StringVar(&configFileType, "configType", "yaml", "yaml/json")
	flag.StringVar(&secret, "secret", "", "Secret key")
	flag.StringVar(&destination, "dest", "dest.yaml", "Destination path")
	flag.Parse()

	if len(secret) == 0 {
		fmt.Println("use --help")
		return
	}

	var unmarshal func(in []byte, out interface{}) (err error)
	var marshal func(v interface{}) ([]byte, error)

	switch configFileType {
	case "json", "JSON":
		unmarshal = json.Unmarshal
		marshal = json.Marshal
	default:
		unmarshal = yaml.Unmarshal
		marshal = yaml.Marshal
	}

	var payload map[string]interface{}
	err := internal.UnmarshalFromFile(configFile, &payload, unmarshal)

	if err != nil {
		return
	}

	result := make(map[string]interface{})
	encryptProperties(secret, payload, &result)

	err = internal.MarshalToFile(destination, result, marshal)

	if err != nil {
		log.Printf("Error in MarshalToFile %v", err)
	}
}

func encryptProperties(secret string, source map[string]interface{}, destination *map[string]interface{}) {
	for key, value := range source {
		fmt.Println(key, value)
		switch v := value.(type) {
		default:
			(*destination)[key], _ = internal.EncryptBytesToBase64([]byte(fmt.Sprint(v)), secret)
		case map[string]interface{}:
			sub := make(map[string]interface{})
			(*destination)[key] = sub
			encryptProperties(secret, value.(map[string]interface{}), &sub)
		}
	}
}
