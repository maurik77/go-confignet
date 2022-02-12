package extensions

// IConfigurationDecrypter is the configuration decrypter interface
type IConfigurationDecrypter interface {
	Init(configurationBuilder IConfigurationBuilder)
	Decrypt(encryptedValue ...string) (decryptedValue string, err error)
}
