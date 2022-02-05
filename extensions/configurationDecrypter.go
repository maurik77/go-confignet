package extensions

// IConfigurationDecrypter is the configuration decrypter interface
type IConfigurationDecrypter interface {
	Decrypt(encryptedValue ...string) (decryptedValue string, err error)
}
