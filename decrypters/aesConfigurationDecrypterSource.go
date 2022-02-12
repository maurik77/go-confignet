package decrypters

import "github.com/maurik77/go-confignet/extensions"

const (
	// DecrypterAesIdentifier is the environment variable containing the UniqueIdentifier of the decrypter
	DecrypterAesIdentifier = "aes"
)

// AesConfigurationDecrypterSource is able to create AesConfigurationDecrypter starting from the provider settings
type AesConfigurationDecrypterSource struct {
}

// NewConfigurationDecrypter creates AesConfigurationDecrypter starting from the provider settings
func (decrypterSource *AesConfigurationDecrypterSource) NewConfigurationDecrypter(settings extensions.DecrypterSettings) extensions.IConfigurationDecrypter {
	if settings.Name != decrypterSource.GetUniqueIdentifier() {
		panic("AesConfigurationDecrypterSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + decrypterSource.GetUniqueIdentifier())
	}

	secret := settings.GetPropertyValue("secret", "").(string)

	return &AesConfigurationDecrypter{
		Secret: secret,
	}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (decrypterSource *AesConfigurationDecrypterSource) GetUniqueIdentifier() string {
	return DecrypterAesIdentifier
}
