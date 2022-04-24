package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

const SECRET = "_____TO BE REPLACED_____"

func Encrypt(value string) string {
	block, err := aes.NewCipher([]byte(SECRET))
	if err != nil {
		panic(fmt.Errorf("could not create new cipher: %w", err))
	}

	bValue := []byte(value)
	cipherValue := make([]byte, aes.BlockSize+len(bValue))
	iv := cipherValue[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		panic(fmt.Errorf("could not read random number: %w", err))
	}

	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(cipherValue[aes.BlockSize:], bValue)

	return string(cipherValue)
}

func Decrypt(value string) string {
	block, err := aes.NewCipher([]byte(SECRET))
	if err != nil {
		panic(fmt.Errorf("could not create new cipher: %w", err))
	}

	bValue := []byte(value)
	if len(bValue) < aes.BlockSize {
		panic(fmt.Errorf("could not decrypt value: [%s ]. It is too short", value))
	}

	iv := bValue[:aes.BlockSize]
	bValue = bValue[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(bValue, bValue)

	return string(bValue)
}
