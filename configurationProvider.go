package confignet

// IConfigurationProvider is configuration provider interface
type IConfigurationProvider interface {
	Load()
	GetData() map[string]string
	GetSeparator() string
}

// IChainedConfigurationProvider is configuration provider interface
type IChainedConfigurationProvider interface {
	IConfigurationProvider
	Add(source IConfigurationProvider)
}
