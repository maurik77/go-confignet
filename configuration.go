package confignet

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
)

type InvalidBindError struct {
	Type reflect.Type
}

var (
	conf extensions.IConfiguration
	lock = &sync.Mutex{}
)

func (e *InvalidBindError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

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

// Bind applies the configuration to the given object using the default configuration providers (AddDefaultConfigurationProviders)
func Bind(section string, target interface{}) error {
	initializeDefaultConfig()
	return conf.Bind(section, target)
}

func initializeDefaultConfig() {
	if conf == nil {
		lock.Lock()
		defer lock.Unlock()
		if conf == nil {
			var confBuilder extensions.IConfigurationBuilder = &ConfigurationBuilder{}
			confBuilder.AddDefaultConfigurationProviders()
			conf = confBuilder.Build()
		}
	}
}

// Bind applies the configuration to the given object
func (conf *Configuration) Bind(section string, target interface{}) error {
	rv := reflect.ValueOf(target)

	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidBindError{reflect.TypeOf(target)}
	}

	for _, p := range conf.configurationProvidersInfo {
		props := filterProperties(section, p)
		conf.bindProps(p, props, target)
	}

	return nil
}

func (conf *Configuration) bindProps(configInfo extensions.ConfigurationProviderInfo, props map[string]string, target interface{}) {
	for key, value := range props {
		parts := strings.Split(key, configInfo.Provider.GetSeparator())
		internal.Unmarshal(target, value, parts...)
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
