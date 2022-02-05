package decrypters

import "github.com/Maurik77/go-confignet/internal"

// AesConfigurationDecrypter .
type AesConfigurationDecrypter struct {
	Secret string
}

// Decrypt decrypts the input encrypted string using aes256 algorithm
func (decrypter *AesConfigurationDecrypter) Decrypt(encryptedValue ...string) (decryptedValue string, err error) {
	decryptedBytes, err := internal.DecryptBase64ToBytes(encryptedValue[0], decrypter.Secret)

	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
