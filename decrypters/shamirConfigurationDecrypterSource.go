package decrypters

import "github.com/Maurik77/go-confignet/extensions"

// ShamirConfigurationDecrypterSource is able to create ShamirConfigurationDecrypter starting from the provider settings
type ShamirConfigurationDecrypterSource struct {
}

// NewConfigurationDecrypter creates AesConfigurationDecrypter starting from the provider settings
func (decrypterSource *ShamirConfigurationDecrypterSource) NewConfigurationDecrypter(settings extensions.DecrypterSettings) extensions.IConfigurationDecrypter {
	if settings.Name != decrypterSource.GetUniqueIdentifier() {
		panic("ShamirConfigurationDecrypterSource: settings of configuration source " + settings.Name + " has been passed to the configuration source with unique identifier " + decrypterSource.GetUniqueIdentifier())
	}

	return &ShamirConfigurationDecrypter{}
}

// GetUniqueIdentifier returns the unique identifier of the configuration provider source. It will be use in the settings file
func (decrypterSource *ShamirConfigurationDecrypterSource) GetUniqueIdentifier() string {
	return "shamir"
}
