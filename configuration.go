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
		*v, _ = strconv.Atoi(value)
	case *int8:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			*v = int8(parsed)
		}
	case *int16:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			*v = int16(parsed)
		}
	case *int64:
		*v, _ = strconv.ParseInt(value, 10, 64)
	case *uint:
		parsed, err := strconv.ParseUint(value, 10, 0)
		if err == nil {
			*v = uint(parsed)
		}
	case *uint8:
		parsed, err := strconv.ParseUint(value, 10, 16)
		if err == nil {
			*v = uint8(parsed)
		}
	case *uint16:
		parsed, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			*v = uint16(parsed)
		}
	case *uint64:
		*v, _ = strconv.ParseUint(value, 10, 64)
	case *bool:
		*v, _ = strconv.ParseBool(value)
	case *time.Time:
		*v, _ = time.Parse(time.RFC3339Nano, value)
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
