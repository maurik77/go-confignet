package confignet

import (
	extensions "confignet/extensions"
	providers "confignet/providers"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	for _, providerSettings := range settings.Providers {
		if configurationSource, ok := configurationSources[providerSettings.Name]; ok {
			conf.Add(configurationSource.NewConfigurationProvider(providerSettings))
		} else {
			log.Printf("ConfigurationBuilder: unable to find configuration source with unique identifier '%v'", providerSettings.Name)
		}
	}
}

// ConfigureConfigurationProviders adds the configuration providers reading from settings.json or settings.yaml file
func (conf *ConfigurationBuilder) ConfigureConfigurationProviders() {
	conf.ConfigureConfigurationProvidersFromJSONConfig("")
	conf.ConfigureConfigurationProvidersFromYamlConfig("")
}

// ConfigureConfigurationProvidersFromJSONConfig adds the configuration providers reading from settings.json file
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromJSONConfig(jsonPath string) {
	if jsonPath == "" {
		jsonPath = "settings.json"
	}

	if _, err := os.Stat(jsonPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("ConfigurationBuilder:File not found %v", jsonPath)
		return
	}

	content, err := ioutil.ReadFile(jsonPath)

	if err != nil {
		log.Printf("ConfigurationBuilder:Error when opening file '%v': '%v'", jsonPath, err)
		return
	}

	var settings extensions.Settings
	err = json.Unmarshal(content, &settings)
	if err != nil {
		log.Println("ConfigurationBuilder:Error during Unmarshal(): ", err)
	}

	conf.ConfigureConfigurationProvidersFromSettings(settings)
}

// ConfigureConfigurationProvidersFromJsonConfig adds the configuration providers reading from settings.yaml file
func (conf *ConfigurationBuilder) ConfigureConfigurationProvidersFromYamlConfig(jsonPath string) {
	if jsonPath == "" {
		jsonPath = "settings.yaml"
	}

	if _, err := os.Stat(jsonPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("ConfigurationBuilder:File not found %v", jsonPath)
		return
	}

	content, err := ioutil.ReadFile(jsonPath)

	if err != nil {
		log.Printf("ConfigurationBuilder:Error when opening file '%v': '%v'", jsonPath, err)
		return
	}

	var settings extensions.Settings
	err = json.Unmarshal(content, &settings)
	if err != nil {
		log.Println("ConfigurationBuilder:Error during Unmarshal(): ", err)
	}

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
