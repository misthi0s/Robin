package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func encryptConfig(config []byte, keyString string) string {
	h := sha256.New()
	h.Write([]byte(keyString))
	block, err := aes.NewCipher(h.Sum(nil))
	if err != nil {
		fmt.Printf("[*] Error creating AES cipher block. Error message: %s\n", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, config, nil)

	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded
}
