package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

const SECRET = "_____TO BE REPLACED_____"

func Encrypt(value string) string {
	block, err := aes.NewCipher([]byte(SECRET))
	if err != nil {
		panic(err)
	}

	bValue := []byte(value)
	cipherValue := make([]byte, aes.BlockSize+len(bValue))
	iv := cipherValue[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		panic(err)
	}

	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(cipherValue[aes.BlockSize:], bValue)

	return string(cipherValue)
}

func Decrypt(value string) string {
	block, err := aes.NewCipher([]byte(SECRET))
	if err != nil {
		panic(err)
	}

	bValue := []byte(value)
	if len(bValue) < aes.BlockSize {
		panic("Could not decrypt value: [" + value + "]. It is too short")
	}

	iv := bValue[:aes.BlockSize]
	bValue = bValue[aes.BlockSize:]

	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(bValue, bValue)

	return string(bValue)
}
