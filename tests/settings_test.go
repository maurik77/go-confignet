package tests

import (
	"confignet"
	extensions "confignet/extensions"
	"fmt"

	"testing"
)

func TestConfigureConfigurationProviders(t *testing.T) {
	var confBuilder extensions.IConfigurationBuilder = &confignet.ConfigurationBuilder{}

	confBuilder.ConfigureConfigurationProviders()

	config := confBuilder.Build()

	fmt.Sprintln(config.GetProviders())
}
