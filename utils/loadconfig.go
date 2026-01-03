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

// Config represents the structure of endmi.json
type Config struct {
	TempDir string `json:"TempDir"`
}

// getHomeDir resolves the user's home directory.
// Returns an error instead of silently falling back.
func getHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		return "", errors.New("unable to resolve user home directory")
	}
	return homeDir, nil
}

// GetConfigDir returns the full path to ~/.endmi
func GetConfigDir() (string, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configDirName), nil
}

// GetConfigFilePath returns the full path to ~/.endmi/endmi.json
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

// CreateConfigPathIfNotExists ensures ~/.endmi exists
func CreateConfigPathIfNotExists() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0755)
}

// GenerateDefaultConfig builds the default configuration dynamically
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

// WriteConfig writes the default configuration file.
// Fails if the file already exists.
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

// CheckConfigExists returns true if ~/.endmi/endmi.json exists
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

// EnsureConfig ensures the config directory and file exist.
// Safe to call multiple times.
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
