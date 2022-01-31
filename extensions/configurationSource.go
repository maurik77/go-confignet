package extensions

// IConfigurationSource is the interface of the configuration source
type IConfigurationSource interface {
	GetUniqueIdentifier() string
	NewConfigurationProvider(settings ProviderSettings) IConfigurationProvider
}

type Settings struct {
	Providers        []ProviderSettings
	ChainedProviders []ChainedProviderSettings
}

type ChainedProviderSettings struct {
	ProviderSettings
	Providers []ProviderSettings
}

type ProviderSettings struct {
	Name       string
	Properties map[string]interface{}
}

func (providerSettings *ProviderSettings) GetPropertyValue(key string, defaultValue interface{}) interface{} {
	if value, ok := providerSettings.Properties[key]; ok {
		return value
	} else {
		return defaultValue
	}
}
