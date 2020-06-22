package aes256cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

func Decode(encrypted string, password string) (string, error) {

	key := []byte(password)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, err = unpad(cipherText, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", cipherText), nil
}

func unpad(data []byte, blockSize uint) ([]byte, error) {
	if blockSize < 1 {
		return nil, fmt.Errorf("block size looks wrong")
	}

	if uint(len(data))%blockSize != 0 {
		return nil, fmt.Errorf("data isn't aligned to blockSize")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	paddingLength := int(data[len(data)-1])

	if (len(data) - paddingLength) < 0 {
		return nil, fmt.Errorf("error padding or data length")
	}

	for _, el := range data[len(data)-paddingLength:] {
		if el != byte(paddingLength) {
			return nil, fmt.Errorf("padding had malformed entries. Have '%x', expected '%x'", paddingLength, el)
		}
	}

	return data[:len(data)-paddingLength], nil
}
