package decrypters

import (
	"encoding/base64"

	"github.com/lafriks/go-shamir"
)

// ShamirConfigurationDecrypter .
type ShamirConfigurationDecrypter struct {
}

// Decrypt decrypts the input encrypted string using aes256 algorithm
func (decrypter *ShamirConfigurationDecrypter) Decrypt(encryptedValue ...string) (decryptedValue string, err error) {

	parts := [][]byte{}

	for _, value := range encryptedValue {
		decodedString, err := base64.StdEncoding.DecodeString(value)

		if err != nil {
			continue
		}

		parts = append(parts, decodedString)
	}

	decryptedBytes, err := shamir.Combine(parts...)

	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
