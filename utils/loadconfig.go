// loadconfig.go
// Configuration management for endmi settings stored in ~/.endmi/endmi.json.
//
// Types:
//   - Config: represents the structure of endmi.json
//
// Functions:
//   - getHomeDir: resolves the user's home directory
//   - GetConfigDir: returns the full path to ~/.endmi
//   - GetConfigFilePath: returns the full path to ~/.endmi/endmi.json
//   - CreateConfigPathIfNotExists: ensures ~/.endmi exists
//   - GenerateDefaultConfig: builds the default configuration dynamically
//   - WriteConfig: writes the default configuration file
//   - CheckConfigExists: returns true if ~/.endmi/endmi.json exists
//   - EnsureConfig: ensures the config directory and file exist
package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	configDirName  = ".endmi"
	configFileName = "endmi.json"
)

type Config struct {
	TempDir string `json:"TempDir"`
}

func getHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		return "", errors.New("unable to resolve user home directory")
	}
	return homeDir, nil
}

func GetConfigDir() (string, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configDirName), nil
}

func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

func CreateConfigPathIfNotExists() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0755)
}

func GenerateDefaultConfig() ([]byte, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	tempDir := filepath.Join(configDir, "tmp")

	cfg := Config{
		TempDir: tempDir,
	}

	return json.MarshalIndent(cfg, "", "\t")
}

func WriteConfig() error {
	if err := CreateConfigPathIfNotExists(); err != nil {
		return err
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := GenerateDefaultConfig()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		configPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	return os.WriteFile(configPath, data, 0644)
}

func CheckConfigExists() (bool, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func EnsureConfig() error {
	exists, err := CheckConfigExists()
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	if err := WriteConfig(); err != nil {
		return err
	}

	// Ensure TempDir exists
	data, err := GenerateDefaultConfig()
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return os.MkdirAll(cfg.TempDir, 0755)
}
