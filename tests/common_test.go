package tests

import (
	"confignet"
	"fmt"
	"testing"
)

func TestConfigurationProviderSources(t *testing.T) {
	configSources := confignet.ConfigurationSources()

	expectedValues := []string{
		"cmdline",
		"env",
		"json",
		"yaml",
		"keyvault",
	}

	for _, tc := range expectedValues {
		t.Run(fmt.Sprintf("Expected value=%v", tc), func(t *testing.T) {
			if !sliceContains(configSources, tc) {
				t.Fatalf("Unable to find expected value %v", tc)
			} else {
				t.Logf("Success !")
			}
		})
	}
}

func sliceContains(slice []string, value string) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	return false
}
