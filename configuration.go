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
	// log.Printf("fillObject -> parts: %v, value: %v", parts, value)

	if parent.Kind() != reflect.Struct {
		log.Printf("Configuration:Parent is not a valid struct %v. Parts %v", parent, parts)
		return
	}

	fieldName := parts[0]
	nestedField := parent.FieldByName(fieldName)

	if !nestedField.IsValid() {
		log.Printf("Configuration:Unable to find field %v in the object %v", fieldName, nestedField)
		return
	}

	if nestedField.Kind() == reflect.Ptr {
		if nestedField.IsNil() {
			defaultValue := reflect.New(nestedField.Type().Elem()).Elem()
			nestedField.Set(defaultValue.Addr())
		}
		nestedField = nestedField.Elem()
	}

	switch {
	case len(parts) == 1: // property
		fillField(nestedField, value, 0)
	case nestedField.Kind() == reflect.Slice:
		conf.fillSlice(configInfo, fieldName, nestedField, value, parts...)
	case nestedField.Kind() == reflect.Array:
		conf.fillArray(configInfo, fieldName, nestedField, value, parts...)
	// case nestedField.Kind() == reflect.Map:
	// 	log.Printf("fillObject::Map: parts %v value %v", parts, value)
	default: // nested object
		conf.fillObject(configInfo, nestedField, value, parts[1:]...)
	}
}

func (conf *Configuration) fillArray(configInfo extensions.ConfigurationProviderInfo, fieldName string, nestedField reflect.Value, value string, parts ...string) {
	// log.Printf("fillArray -> fieldName: %v, parts: %v, value: %v", fieldName, parts, value)

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Configuration:Unable to parse index %v for field %v in the object %v", parts[1], fieldName, nestedField)
		return
	}

	if index >= nestedField.Len() {
		log.Printf("Configuration:Unable to assign value %v to the field %v with index %v because is out of range. (Array length  %v)", value, fieldName, index, nestedField.Len())
		return
	}

	if len(parts) == 2 {
		fillField(nestedField, value, index)
	} else {
		conf.fillObject(configInfo, nestedField.Index(index), value, parts[2:]...)
	}
}

func (conf *Configuration) fillSlice(configInfo extensions.ConfigurationProviderInfo, fieldName string, nestedField reflect.Value, value string, parts ...string) {
	// log.Printf("fillSlice -> fieldName: %v, parts: %v, value: %v", fieldName, parts, value)

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Printf("Configuration:Unable to parse index %v for field %v in the object %v", parts[1], fieldName, nestedField)
		return
	}

	if nestedField.IsNil() {
		elemType := nestedField.Type().Elem()
		elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
		nestedField.Set(elemSlice)
	}

	if index >= nestedField.Len() {
		newElements := make([]reflect.Value, index+1-nestedField.Len())
		for i := range newElements {
			newElements[i] = reflect.New(nestedField.Type().Elem()).Elem()
		}
		nestedField.Set(reflect.Append(nestedField, newElements...))
	}

	if len(parts) == 2 {
		fillField(nestedField, value, index)
	} else {
		conf.fillObject(configInfo, nestedField.Index(index), value, parts[2:]...)
	}
}

func fillField(field reflect.Value, value string, index int) {
	// log.Printf("fillField -> index: %v, value: %v", index, value)

	switch field.Kind() {
	case reflect.Array:
		item := field.Index(index)
		fillField(item, value, -1)
		return
	case reflect.Slice:
		item := field.Index(index)
		fillField(item, value, -1)
		return
	case reflect.Ptr:
		if field.IsNil() {
			defaultValue := reflect.New(field.Type().Elem()).Elem()
			field.Set(defaultValue.Addr())
		}
		fillField(field.Elem(), value, -1)
	default:
		switch field.Interface().(type) {
		case string:
			field.SetString(value)
		case int, int8, int16, int32, int64:
			if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
				field.SetInt(intValue)
			} else {
				log.Printf("Configuration:Unable to parse Int the value %v", value)
			}
		case uint, uint8, uint16, uint32, uint64:
			if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
				field.SetUint(uintValue)
			} else {
				log.Printf("Configuration:Unable to parse Uint the value %v", value)
			}
		case bool:
			if boolValue, err := strconv.ParseBool(value); err == nil {
				field.SetBool(boolValue)
			} else {
				log.Printf("Configuration:Unable to parse Bool the value %v", value)
			}
		case time.Time:
			if timeValue, err := time.Parse(time.RFC3339Nano, value); err == nil {
				field.Set(reflect.ValueOf(timeValue))
			} else {
				log.Printf("Configuration:Unable to parse Time (format: %v) the value %v", time.RFC3339Nano, value)
			}
		}
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
