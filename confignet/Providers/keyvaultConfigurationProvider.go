package providers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// KeyvaultConfigurationProvider loads configuration from Azure Key Vault
type KeyvaultConfigurationProvider struct {
	Prefix       string
	RemovePrefix bool
	TenantID     string
	ClientID     string
	ClientSecret string
	BaseURL      string
	data         map[string]string
}

// Load configuration from Azure Key Vault
func (provider *KeyvaultConfigurationProvider) Load() {
	provider.data = make(map[string]string)

	cred, err := provider.getCredential()

	if err != nil {
		log.Println("KeyvaultConfigurationProvider:Unable to retrieve the token with the provided credentials")
	}

	client, err := azsecrets.NewClient(provider.BaseURL, cred, nil)

	if err != nil {
		log.Println("KeyvaultConfigurationProvider:Unable to connect to keyvault with the provided credentials and base url", provider.BaseURL)
	}

	pager := client.ListSecrets(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		for _, secret := range resp.Secrets {
			key := strings.TrimPrefix(*secret.ID, fmt.Sprintf("%v/secrets/", provider.BaseURL))

			if provider.Prefix != "" && !strings.HasPrefix(key, provider.Prefix) {
				continue
			}

			if provider.Prefix != "" && provider.RemovePrefix {
				key = strings.TrimPrefix(key, provider.Prefix)
			}

			resp, err := client.GetSecret(context.Background(), key, nil)
			if err == nil {
				provider.data[key] = *resp.Value
			}
		}
	}
}

// GetData provides the loaded data
func (provider *KeyvaultConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *KeyvaultConfigurationProvider) GetSeparator() string {
	return "--"
}

func (provider *KeyvaultConfigurationProvider) getCredential() (azcore.TokenCredential, error) {
	if provider.ClientID != "" && provider.TenantID != "" && provider.ClientSecret != "" && provider.BaseURL != "" {
		return azidentity.NewClientSecretCredential(provider.TenantID, provider.ClientID, provider.ClientSecret, nil)
	}

	return azidentity.NewDefaultAzureCredential(nil)
}
