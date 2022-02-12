package decrypters

import (
	"github.com/maurik77/go-confignet/extensions"
	"github.com/maurik77/go-confignet/internal"
)

// AesConfigurationDecrypter .
type AesConfigurationDecrypter struct {
	Secret           string
	ConfigFileType   string
	ConfigFilePath   string
	SecretConfigPath string
}

// Init the decrypter
func (decrypter *AesConfigurationDecrypter) Init(configurationBuilder extensions.IConfigurationBuilder) {
	if len(decrypter.ConfigFileType) > 0 &&
		len(decrypter.ConfigFilePath) > 0 &&
		len(decrypter.SecretConfigPath) > 0 {

		configurationBuilder.ConfigureConfigurationProviders(decrypter.ConfigFileType, decrypter.ConfigFilePath)
		config := configurationBuilder.Build()
		secret := config.GetValue(decrypter.SecretConfigPath)

		if len(secret) > 0 {
			decrypter.Secret = secret
		}
	}
}

// Decrypt decrypts the input encrypted string using aes256 algorithm
func (decrypter *AesConfigurationDecrypter) Decrypt(encryptedValue ...string) (decryptedValue string, err error) {
	decryptedBytes, err := internal.DecryptBase64ToBytes(encryptedValue[0], decrypter.Secret)

	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
