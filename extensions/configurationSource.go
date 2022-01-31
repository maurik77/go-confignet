package extensions

// IConfigurationSource is the interface of the configuration source
type IConfigurationSource interface {
	GetUniqueIdentifier() string
	NewConfigurationProvider(settings ProviderSettings) IConfigurationProvider
}

// Settings contains information usefull to configure the configuration providers
type Settings struct {
	Providers []ProviderSettings `yaml:"providers" json:"providers"`
}

// ProviderSettings contains information usefull to configure a specific configuration provider
type ProviderSettings struct {
	Name       string                 `yaml:"name" json:"name"`
	Properties map[string]interface{} `yaml:"properties" json:"properties"`
	Providers  []ProviderSettings     `yaml:"providers" json:"providers"`
}

// GetPropertyValue return the found value or the default if the key doesn't exist in the collection
func (providerSettings *ProviderSettings) GetPropertyValue(key string, defaultValue interface{}) interface{} {
	if value, ok := providerSettings.Properties[key]; ok {
		return value
	}

	return defaultValue
}
