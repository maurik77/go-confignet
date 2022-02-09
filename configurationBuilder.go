package confignet

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Maurik77/go-confignet/extensions"
	"github.com/Maurik77/go-confignet/internal"
	"github.com/Maurik77/go-confignet/providers"
	"gopkg.in/yaml.v2"
)

const (
	// EnvConfigFileType is the environment variable containing the type of the settings file: yaml or json
	EnvConfigFileType = "confignet_configfiletype"
	// EnvConfigFilePath is the environment variable containing the file path
	EnvConfigFilePath = "confignet_configfilepath"
)

// ConfigurationBuilder is the concrete implementation
type ConfigurationBuilder struct {
	configurationProvidersInfo []extensions.ConfigurationProviderInfo
}

// Add adds the configuration provider to the inner collection
func (conf *ConfigurationBuilder) Add(source extensions.IConfigurationProvider) {
	conf.configurationProvidersInfo = append(conf.configurationProvidersInfo, extensions.ConfigurationProviderInfo{Provider: source})
	log.Printf("ConfigurationBuilder:Added configuration provider '%T', Separator:'%v'\n", source, source.GetSeparator())
}

// AddWithEncrypter adds the configuration provider and the decrypter to the inner collection
func (conf *ConfigurationBuilder) AddWithEncrypter(source extensions.IConfigurationProvider, decrypter extensions.IConfigurationDecrypter) {
	conf.configurationProvidersInfo = append(conf.configurationProvidersInfo, extensions.ConfigurationProviderInfo{Provider: source, Decrypter: decrypter})
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

	configureConfigurationProvidersFromSettings(settings.Providers, conf)
}

func configureConfigurationProvidersFromSettings(settings []extensions.ProviderSettings, configurationProvidersCollection extensions.IConfigurationProviderCollection) {
	for _, providerSettings := range settings {
		if configurationSource, ok := configurationSources[providerSettings.Name]; ok {
			provider := configurationSource.NewConfigurationProvider(providerSettings)

			if decrypterSource, ok := decrypterSources[providerSettings.Decrypter.Name]; ok {
				configurationProvidersCollection.AddWithEncrypter(provider, decrypterSource.NewConfigurationDecrypter(providerSettings.Decrypter))
			} else {
				configurationProvidersCollection.Add(provider)
			}

			if chainedConfigurationProvider, ok := provider.(extensions.IConfigurationProviderCollection); ok {
				configureConfigurationProvidersFromSettings(providerSettings.Providers, chainedConfigurationProvider)
			}

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
	case "json", "JSON":
		conf.ConfigureConfigurationProvidersFromJSONConfig(configFilePath)
	default:
		conf.ConfigureConfigurationProvidersFromYamlConfig(configFilePath)
	}
}

// ConfigureConfigurationProvidersFromJSONConfig adds the configuration providers reading from settings.json file
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromJSONConfig(jsonPath string) {
	if len(jsonPath) == 0 {
		jsonPath = "settings.json"
	}
	var settings extensions.Settings
	internal.UnmarshalFromFile(jsonPath, &settings, json.Unmarshal)
	conf.ConfigureConfigurationProvidersFromSettings(settings)
}

// ConfigureConfigurationProvidersFromYamlConfig adds the configuration providers reading from settings.yaml file
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromYamlConfig(yamlPath string) {
	if len(yamlPath) == 0 {
		yamlPath = "settings.yaml"
	}
	var settings extensions.Settings
	internal.UnmarshalFromFile(yamlPath, &settings, yaml.Unmarshal)
	conf.ConfigureConfigurationProvidersFromSettings(settings)
}

// Build invokes the load function of each configuration provider and return the Configuration object
func (conf *ConfigurationBuilder) Build() extensions.IConfiguration {
	for _, confProvider := range conf.configurationProvidersInfo {
		confProvider.Provider.Load(confProvider.Decrypter)
	}

	result := Configuration{
		configurationProvidersInfo: conf.configurationProvidersInfo,
	}

	return &result
}
