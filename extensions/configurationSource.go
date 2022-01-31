package extensions

// IConfigurationSource is the interface of the configuration source
type IConfigurationSource interface {
	GetUniqueIdentifier() string
	NewConfigurationProvider(settings ProviderSettings) IConfigurationProvider
}

type Settings struct {
	Providers        []ProviderSettings        `yaml:"providers" json:"providers"`
	ChainedProviders []ChainedProviderSettings `yaml:"chainedProviders" json:"chainedProviders"`
}

type ChainedProviderSettings struct {
	ProviderSettings `yaml:",inline"`
	Providers        []ProviderSettings `yaml:"providers" json:"providers"`
}

type ProviderSettings struct {
	Name       string                 `yaml:"name" json:"name"`
	Properties map[string]interface{} `yaml:"properties" json:"properties"`
}

func (providerSettings *ProviderSettings) GetPropertyValue(key string, defaultValue interface{}) interface{} {
	if value, ok := providerSettings.Properties[key]; ok {
		return value
	} else {
		return defaultValue
	}
}