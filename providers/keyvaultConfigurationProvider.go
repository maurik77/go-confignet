package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/maurik77/go-confignet/extensions"
	"github.com/rs/zerolog/log"
)

// KeyVaultConfigurationProvider loads configuration from Azure Key Vault
type KeyVaultConfigurationProvider struct {
	Prefix       string
	RemovePrefix bool
	TenantID     string
	ClientID     string
	ClientSecret string
	BaseURL      string
	data         map[string]string
}

// Load configuration from Azure Key Vault
func (provider *KeyVaultConfigurationProvider) Load(decrypter extensions.IConfigurationDecrypter) {
	provider.data = make(map[string]string)

	cred, err := provider.getCredential()

	if err != nil {
		log.Err(err).Msg("KeyVaultConfigurationProvider:Unable to retrieve the token with the provided credentials")
	}

	client, err := azsecrets.NewClient(provider.BaseURL, cred, nil)

	if err != nil {
		log.Err(err).Msgf("KeyVaultConfigurationProvider:Unable to connect to keyvault with the provided credentials and base url: %v", provider.BaseURL)
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
			if err != nil {
				log.Err(err).Msgf("KeyVaultConfigurationProvider:Error retrieving key %v", key)
				continue
			}

			value := *resp.Value

			if decrypter != nil {
				var err error
				value, err = decrypter.Decrypt(value)

				if err != nil {
					log.Err(err).Msgf("KeyVaultConfigurationProvider:Error calling decryption for key %v", key)
				}
			}

			provider.data[key] = value
		}
	}
}

// GetData provides the loaded data
func (provider *KeyVaultConfigurationProvider) GetData() map[string]string {
	return provider.data
}

// GetSeparator provides the separator that it uses to store nested object
func (provider *KeyVaultConfigurationProvider) GetSeparator() string {
	return "--"
}

func (provider *KeyVaultConfigurationProvider) getCredential() (azcore.TokenCredential, error) {
	if provider.ClientID != "" && provider.TenantID != "" && provider.ClientSecret != "" && provider.BaseURL != "" {
		return azidentity.NewClientSecretCredential(provider.TenantID, provider.ClientID, provider.ClientSecret, nil)
	}

	return azidentity.NewDefaultAzureCredential(nil)
}
