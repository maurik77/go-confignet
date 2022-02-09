package extensions

// IConfigurationProvider is configuration provider interface
type IConfigurationProvider interface {
	Load(decrypter IConfigurationDecrypter)
	GetData() map[string]string
	GetSeparator() string
}

// IChainedConfigurationProvider is configuration provider interface
type IChainedConfigurationProvider interface {
	IConfigurationProvider
	IConfigurationProviderCollection
}
