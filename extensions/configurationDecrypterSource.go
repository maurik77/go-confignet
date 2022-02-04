package extensions

// IConfigurationDecrypterSource is the configuration decrypter source interface
type IConfigurationDecrypterSource interface {
	GetUniqueIdentifier() string
	NewConfigurationDecrypter(settings DecrypterSettings) IConfigurationDecrypter
}

// DecrypterSettings contains information usefull to configure a specific configuration decrypter
type DecrypterSettings struct {
	Name       string                 `yaml:"name" json:"name"`
	Properties map[string]interface{} `yaml:"properties" json:"properties"`
}

// GetPropertyValue return the found value or the default if the key doesn't exist in the collection
func (decrypterSettings *DecrypterSettings) GetPropertyValue(key string, defaultValue interface{}) interface{} {
	if value, ok := decrypterSettings.Properties[key]; ok {
		return value
	}

	return defaultValue
}
