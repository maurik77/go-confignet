package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

// EncryptBytesToBase64 encrypts string with secret
func EncryptBytesToBase64(src []byte, keyString string) (string, error) {
	data, err := aesEncrypt(src, keyString)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b)
	_, err = encoder.Write(data)
	encoder.Close()

	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func aesEncrypt(src []byte, keyString string) ([]byte, error) {
	key := createHash(keyString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aesDecrypter:[Error] %s", err)
	}

	if src == nil {
		return nil, fmt.Errorf("aesDecrypter:[Error] plain content empty")
	}
	initialVector := createHash(string(src))[0:16]

	ecb := cipher.NewCBCEncrypter(block, initialVector)
	src = pkcs5Padding(src, block.BlockSize())
	encryptedValue := make([]byte, len(src))
	ecb.CryptBlocks(encryptedValue, src)
	return append(initialVector, encryptedValue...), nil
}

// DecryptBase64ToBytes decrypts string with secret
func DecryptBase64ToBytes(encryptedString string, keyString string) ([]byte, error) {

	data, err := base64ToBytes(encryptedString)
	if err != nil {
		return nil, err
	}

	return aesDecrypt(data, keyString)
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func base64ToBytes(encryptedString string) ([]byte, error) {
	src := []byte(encryptedString)
	r := bytes.NewReader(src)
	input := base64.NewDecoder(base64.StdEncoding, r)
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func aesDecrypt(_crypt []byte, keyString string) ([]byte, error) {
	key := createHash(keyString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aesDecrypter:[Error] %s", err)
	}

	if len(_crypt) == 0 {
		return nil, fmt.Errorf("aesDecrypter:[Error] plain content empty")
	}

	initialVector := _crypt[0:16]
	crypt := _crypt[16:]

	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	if len(decrypted) == 0 {
		return nil, fmt.Errorf("aesDecrypter:[Error] decrypted too short")
	}

	return pkcs5Trimming(decrypted), nil
}

func pkcs5Trimming(encrypt []byte) []byte {

	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
