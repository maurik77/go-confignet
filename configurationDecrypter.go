package confignet

// IConfigurationDecrypter is configuration decrypt interface interface
type IConfigurationDecrypter interface {
	Decrypt(encryptedValue string) (decryptedValue string, err error)
}
