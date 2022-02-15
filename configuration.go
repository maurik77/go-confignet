package confignet

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/maurik77/go-confignet/extensions"
)

const (
	// SectionSeparator is the separator that is used to specify nested sections
	SectionSeparator string = "/"
)

// Configuration is the concrete implementation
type Configuration struct {
	configurationProvidersInfo []extensions.ConfigurationProviderInfo
}

// GetProviders returns the configured configuration providers
func (conf *Configuration) GetProviders() []extensions.ConfigurationProviderInfo {
	return conf.configurationProvidersInfo
}

// GetValue returns the value of the given configuration
func (conf *Configuration) GetValue(section string) string {
	var result string

	for _, p := range conf.configurationProvidersInfo {
		props := filterProperties(section, p)
		if len(props) == 1 {
			for _, value := range props {
				result = value
				break
			}
		}
	}

	return result
}

// Bind applies the configuration to the given object
func (conf *Configuration) Bind(section string, value interface{}) {
	for _, p := range conf.configurationProvidersInfo {
		props := filterProperties(section, p)
		conf.bindProps(p, props, value)
	}
}

func (conf *Configuration) bindProps(configInfo extensions.ConfigurationProviderInfo, props map[string]string, value interface{}) {
	reflectedType := reflect.ValueOf(value).Elem()

	for key, value := range props {
		parts := strings.Split(key, configInfo.Provider.GetSeparator())
		conf.fillObject(configInfo, reflectedType, value, parts...)
	}
}

func (conf *Configuration) fillObject(configInfo extensions.ConfigurationProviderInfo, parent reflect.Value, value string, parts ...string) {
	fieldName := parts[0]
	nestedField := parent
	_, err := strconv.Atoi(fieldName)

	if err != nil {
		nestedField = parent.FieldByName(fieldName)
	}

	if nestedField == (reflect.Value{}) {
		log.Printf("Configuration:Unable to find field %v in the object %v", fieldName, nestedField)
		return
	}

	if len(parts) == 1 { // property
		fillField(nestedField, value)
	} else { // nested object
		conf.fillObject(configInfo, nestedField, value, parts[1:]...)
	}
}

func fillField(field reflect.Value, value string) {
	if field.Kind() == reflect.Slice {
		log.Printf("Configuration:Field is slice %v in the object %v. Subtype %v", field, field.Addr(), reflect.SliceOf(field.Type()))
	}

	valueInt := field.Addr().Interface()
	switch v := valueInt.(type) {
	case *string:
		*v = value
	case *int:
		var err error
		*v, err = strconv.Atoi(value)
		if err != nil {
			log.Printf("Error in parsing int %v, %v", value, err)
		}
	case *int8:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			*v = int8(parsed)
		} else {
			log.Printf("Error in parsing int8 %v, %v", value, err)
		}
	case *int16:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			*v = int16(parsed)
		} else {
			log.Printf("Error in parsing int16 %v, %v", value, err)
		}
	case *int64:
		var err error
		*v, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Error in parsing int64 %v, %v", value, err)
		}
	case *uint:
		parsed, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			*v = uint(parsed)
		} else {
			log.Printf("Error in parsing uint %v, %v", value, err)
		}
	case *uint8:
		parsed, err := strconv.ParseUint(value, 10, 16)
		if err == nil {
			*v = uint8(parsed)
		} else {
			log.Printf("Error in parsing uint8 %v, %v", value, err)
		}
	case *uint16:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			*v = uint16(parsed)
		} else {
			log.Printf("Error in parsing uint16 %v, %v", value, err)
		}
	case *uint64:
		var err error
		*v, err = strconv.ParseUint(value, 10, 64)
		if err != nil {
			log.Printf("Error in parsing uint64 %v, %v", value, err)
		}
	case *bool:
		var err error
		*v, err = strconv.ParseBool(value)
		if err != nil {
			log.Printf("Error in parsing bool %v, %v", value, err)
		}
	case *time.Time:
		var err error
		*v, err = time.Parse(time.RFC3339Nano, value)
		if err != nil {
			log.Printf("Error in parsing time %v, %v", value, err)
		}
	case *[]interface{}:
		fmt.Println(v)
	}
}

func filterProperties(section string, configInfo extensions.ConfigurationProviderInfo) (data map[string]string) {
	data = map[string]string{}
	properties := configInfo.Provider.GetData()
	separator := configInfo.Provider.GetSeparator()

	if strings.Contains(section, SectionSeparator) {
		section = strings.ReplaceAll(section, SectionSeparator, separator)
	}

	for key, value := range properties {
		if strings.HasPrefix(key, section) {
			mapKey := strings.TrimPrefix(key, fmt.Sprintf("%v%v", section, separator))
			data[mapKey] = value
		}
	}

	return data
}
