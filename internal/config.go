package internal

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed templates/*
var templateFiles embed.FS

type Config struct {
	ConfigDir         string
	MainConfig        map[string]interface{}
	BodyInstructions  string
	TitleInstructions string
	BodyExample       string
	TitleExample      string
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "prgen")
	return configDir, nil
}

func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	config := &Config{
		ConfigDir: configDir,
	}

	err = config.ensureConfigExists()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	err = config.loadFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to load config files: %w", err)
	}

	return config, nil
}

func (c *Config) ensureConfigExists() error {
	err := os.MkdirAll(c.ConfigDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFiles := []string{
		"config.json",
		"body_instructions.md",
		"title_instructions.md",
		"body_example.md",
		"title_example.md",
	}

	for _, filename := range configFiles {
		filePath := filepath.Join(c.ConfigDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			err = c.copyTemplateFile(filename, filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Config) copyTemplateFile(templateName, destPath string) error {
	templatePath := "templates/" + templateName
	data, err := templateFiles.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	err = os.WriteFile(destPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", templateName, err)
	}

	fmt.Printf("Created %s at: %s\n", templateName, destPath)
	return nil
}

func (c *Config) loadFiles() error {
	// Load main config JSON
	configPath := filepath.Join(c.ConfigDir, "config.json")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config.json: %w", err)
	}
	err = json.Unmarshal(configData, &c.MainConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config.json: %w", err)
	}

	// Load instructions and examples
	fileLoaders := map[string]*string{
		"body_instructions.md":  &c.BodyInstructions,
		"title_instructions.md": &c.TitleInstructions,
		"body_example.md":       &c.BodyExample,
		"title_example.md":      &c.TitleExample,
	}

	for filename, target := range fileLoaders {
		filePath := filepath.Join(c.ConfigDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filename, err)
		}
		*target = string(data)
	}

	return nil
}

func (c *Config) GetConfigPath() string {
	return filepath.Join(c.ConfigDir, "config.json")
}

func (c *Config) GetBodyInstructionsPath() string {
	return filepath.Join(c.ConfigDir, "body_instructions.md")
}

func (c *Config) GetTitleInstructionsPath() string {
	return filepath.Join(c.ConfigDir, "title_instructions.md")
}

func (c *Config) GetBodyExamplePath() string {
	return filepath.Join(c.ConfigDir, "body_example.md")
}

func (c *Config) GetTitleExamplePath() string {
	return filepath.Join(c.ConfigDir, "title_example.md")
}
