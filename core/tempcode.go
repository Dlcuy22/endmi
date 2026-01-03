package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dlcuy22/endmi/extensions"
	"github.com/dlcuy22/endmi/utils"
)

// TempCodeManager handles temporary code workspace operations
type TempCodeManager struct {
	App *App
}

// TempProjectMetadata stores metadata about a temporary project
type TempProjectMetadata struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Template  string    `json:"template"`
	Path      string    `json:"path"`
}

// loadConfig loads the endmi configuration
func loadConfig() (*utils.Config, error) {
	configPath, err := utils.GetConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg utils.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// GetTempDir returns the configured temporary directory path
func (tcm *TempCodeManager) GetTempDir() (string, error) {
	cfg, err := loadConfig()
	if err != nil {
		return "", err
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(cfg.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return cfg.TempDir, nil
}

// CreateTempProject creates a new temporary project in the temp workspace
func (tcm *TempCodeManager) CreateTempProject(template extensions.Template, projectName string) (string, error) {
	tempDir, err := tcm.GetTempDir()
	if err != nil {
		return "", err
	}

	// Generate unique project name if not provided
	if projectName == "" {
		projectName = fmt.Sprintf("temp_%d", time.Now().Unix())
	}

	projectPath := filepath.Join(tempDir, projectName)

	// Check if project already exists
	if _, err := os.Stat(projectPath); err == nil {
		return "", fmt.Errorf("temp project '%s' already exists", projectName)
	}

	// Create project using the App's CreateProject method
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create project directory: %w", err)
	}

	// Store original directory since CreateProject expects relative path
	originalApp := *tcm.App
	projectApp := &App{Output: originalApp.Output}

	// Create the project structure
	baseDir := filepath.Join(projectPath, template.RootDir())
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	if err := projectApp.runCommandWithOutput("go", projectPath, "mod", "init", projectName); err != nil {
		return "", err
	}

	for rel, content := range template.Files(projectName) {
		fullPath := filepath.Join(baseDir, rel)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return "", err
		}
	}

	for _, dep := range template.Dependencies() {
		if err := projectApp.runCommandWithOutput("go", projectPath, "get", dep); err != nil {
			return "", err
		}
	}

	if err := projectApp.runCommandWithOutput("go", projectPath, "mod", "tidy"); err != nil {
		return "", err
	}

	// Save metadata
	metadata := TempProjectMetadata{
		Name:      projectName,
		CreatedAt: time.Now(),
		Template:  template.Name(),
		Path:      projectPath,
	}

	if err := tcm.saveMetadata(projectPath, metadata); err != nil {
		// Non-fatal: project is created, metadata is just informational
		fmt.Printf("Warning: failed to save metadata: %v\n", err)
	}

	return projectPath, nil
}

// saveMetadata saves project metadata to a .endmi_meta.json file
func (tcm *TempCodeManager) saveMetadata(projectPath string, metadata TempProjectMetadata) error {
	metaPath := filepath.Join(projectPath, ".endmi_meta.json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metaPath, data, 0644)
}

// loadMetadata loads project metadata from .endmi_meta.json
func (tcm *TempCodeManager) loadMetadata(projectPath string) (*TempProjectMetadata, error) {
	metaPath := filepath.Join(projectPath, ".endmi_meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	var metadata TempProjectMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// ListTempProjects returns a list of all temporary projects
func (tcm *TempCodeManager) ListTempProjects() ([]TempProjectMetadata, error) {
	tempDir, err := tcm.GetTempDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read temp directory: %w", err)
	}

	var projects []TempProjectMetadata
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(tempDir, entry.Name())
		metadata, err := tcm.loadMetadata(projectPath)
		if err != nil {
			// If metadata doesn't exist, create a basic one
			info, err := entry.Info()
			if err != nil {
				continue
			}

			metadata = &TempProjectMetadata{
				Name:      entry.Name(),
				CreatedAt: info.ModTime(),
				Template:  "unknown",
				Path:      projectPath,
			}
		}

		projects = append(projects, *metadata)
	}

	return projects, nil
}

// DeleteTempProject removes a temporary project
func (tcm *TempCodeManager) DeleteTempProject(projectName string) error {
	tempDir, err := tcm.GetTempDir()
	if err != nil {
		return err
	}

	projectPath := filepath.Join(tempDir, projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("temp project '%s' does not exist", projectName)
	}

	return os.RemoveAll(projectPath)
}

// CleanAll removes all temporary projects
func (tcm *TempCodeManager) CleanAll() error {
	tempDir, err := tcm.GetTempDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		projectPath := filepath.Join(tempDir, entry.Name())
		if err := os.RemoveAll(projectPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// PromoteTempProject moves a temporary project to a permanent location
func (tcm *TempCodeManager) PromoteTempProject(projectName, targetPath string) error {
	tempDir, err := tcm.GetTempDir()
	if err != nil {
		return err
	}

	sourcePath := filepath.Join(tempDir, projectName)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("temp project '%s' does not exist", projectName)
	}

	// Check if target already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path '%s' already exists", targetPath)
	}

	// Move the project
	if err := os.Rename(sourcePath, targetPath); err != nil {
		// If rename fails (cross-device), copy and remove
		return fmt.Errorf("failed to move project: %w (consider using copy instead)", err)
	}

	// Remove metadata file from promoted project
	metaPath := filepath.Join(targetPath, ".endmi_meta.json")
	os.Remove(metaPath) // Ignore errors

	return nil
}
