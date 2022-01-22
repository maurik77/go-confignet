package confignet

// IConfigurationProvider is configuration provider interface
type IConfigurationProvider interface {
	Load()
	GetData() map[string]string
	GetSeparator() string
}

// IChaninedConfigurationProvider is configuration provider interface
type IChaninedConfigurationProvider interface {
	IConfigurationProvider
	Add(source IConfigurationProvider)
}
