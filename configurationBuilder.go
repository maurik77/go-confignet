package confignet

import (
	extensions "confignet/extensions"
	providers "confignet/providers"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	EnvConfigFileType = "confignet_configfiletype"
	EnvConfigFilePath = "confignet_configfilepath"
)

// ConfigurationBuilder is the concrete implementation
type ConfigurationBuilder struct {
	configurationProviders []extensions.IConfigurationProvider
	decrypters             map[*extensions.IConfigurationProvider]extensions.IConfigurationDecrypter
}

// Add adds the configuration provider to the inner collection
func (conf *ConfigurationBuilder) Add(source extensions.IConfigurationProvider) {
	conf.configurationProviders = append(conf.configurationProviders, source)
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// AddWithEncrypter adds the configuration provider and the decrypter to the inner collection
func (conf *ConfigurationBuilder) AddWithEncrypter(source extensions.IConfigurationProvider, decrypter extensions.IConfigurationDecrypter) {
	conf.configurationProviders = append(conf.configurationProviders, source)
	conf.decrypters[&source] = decrypter
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n, Decrypter:'%T'", source, source.GetSeparator(), decrypter)
}

// AddDefaultConfigurationProviders adds the default configuration providers
func (conf *ConfigurationBuilder) AddDefaultConfigurationProviders() {
	conf.AddDefaultConfigurationProvidersWithBasePath("")
}

// AddDefaultConfigurationProvidersWithBasePath adds the default configuration providers
func (conf *ConfigurationBuilder) AddDefaultConfigurationProvidersWithBasePath(basePath string) {
	conf.Add(&providers.JSONConfigurationProvider{FilePath: fmt.Sprintf("%v%v", basePath, providers.DefaultJSONFile)})
	conf.Add(&providers.YamlConfigurationProvider{FilePath: fmt.Sprintf("%v%v", basePath, providers.DefaultYAMLFile)})
	conf.Add(&providers.EnvConfigurationProvider{})
	conf.Add(&providers.CmdLineConfigurationProvider{})
	conf.Add(&providers.KeyvaultConfigurationProvider{})
}

// ConfigureConfigurationProvidersFromSettings adds the default configuration providers
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromSettings(settings extensions.Settings) {

	configureConfigurationProvidersFromSettings(settings.Providers, func(providerSettings extensions.ProviderSettings, provider extensions.IConfigurationProvider) {
		conf.Add(provider)
		if chainedConfigurationProvider, ok := provider.(extensions.IChainedConfigurationProvider); ok {
			configureConfigurationProvidersFromSettings(providerSettings.Providers, func(subProviderSettings extensions.ProviderSettings, subprovider extensions.IConfigurationProvider) {
				chainedConfigurationProvider.Add(subprovider)
			})
		}
	})
}

func configureConfigurationProvidersFromSettings(settings []extensions.ProviderSettings, add func(extensions.ProviderSettings, extensions.IConfigurationProvider)) {
	for _, providerSettings := range settings {
		if configurationSource, ok := configurationSources[providerSettings.Name]; ok {
			add(providerSettings, configurationSource.NewConfigurationProvider(providerSettings))
		} else {
			log.Printf("ConfigurationBuilder: unable to find configuration source with unique identifier '%v'", providerSettings.Name)
		}
	}
}

// ConfigureConfigurationProviders adds the configuration providers reading from settings.json or settings.yaml file
func (conf *ConfigurationBuilder) ConfigureConfigurationProviders() {
	var configFileType = os.Getenv(EnvConfigFileType)
	var configFilePath = os.Getenv(EnvConfigFilePath)

	switch configFileType {
	case "JSON":
	case "json":
		conf.ConfigureConfigurationProvidersFromJSONConfig(configFilePath)
	default:
		conf.ConfigureConfigurationProvidersFromYamlConfig(configFilePath)
	}
}

// ConfigureConfigurationProvidersFromJSONConfig adds the configuration providers reading from settings.json file
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromJSONConfig(jsonPath string) {
	settings := unmarshalSettingsFile(jsonPath, "settings.json", yaml.Unmarshal)
	conf.ConfigureConfigurationProvidersFromSettings(settings)
}

// ConfigureConfigurationProvidersFromYamlConfig adds the configuration providers reading from settings.yaml file
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromYamlConfig(yamlPath string) {
	settings := unmarshalSettingsFile(yamlPath, "settings.yaml", yaml.Unmarshal)
	conf.ConfigureConfigurationProvidersFromSettings(settings)
}

// Build invokes the load function of each configuration provider and return the Configuration object
func (conf *ConfigurationBuilder) Build() extensions.IConfiguration {
	for _, confProvider := range conf.configurationProviders {
		confProvider.Load()
	}

	result := Configuration{
		configurationProviders: conf.configurationProviders,
		decrypters:             conf.decrypters,
	}

	return &result
}

func unmarshalSettingsFile(path string, defaultPath string, unmarshal func(in []byte, out interface{}) (err error)) extensions.Settings {
	if path == "" {
		path = "settings.yaml"
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("ConfigurationBuilder:File not found %v", path)
		return extensions.Settings{}
	}

	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Printf("ConfigurationBuilder:Error when opening file '%v': '%v'", path, err)
		return extensions.Settings{}
	}

	var settings extensions.Settings
	err = unmarshal(content, &settings)
	if err != nil {
		log.Println("ConfigurationBuilder:Error during Unmarshal(): ", err)
	}

	return settings
}
