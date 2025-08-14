package internal

import "fmt"

// Provider represents an AI provider interface
type Provider interface {
	GeneratePRContent(config *Config, diff string) (title, body string, err error)
}

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct{}

// GeneratePRContent generates PR content using OpenAI
func (p *OpenAIProvider) GeneratePRContent(config *Config, diff string) (title, body string, err error) {
	return GeneratePRContent(config, diff)
}

// ClaudeProvider implements the Provider interface for Claude Code CLI
type ClaudeProvider struct{}

// GeneratePRContent generates PR content using Claude Code CLI
func (p *ClaudeProvider) GeneratePRContent(config *Config, diff string) (title, body string, err error) {
	return GeneratePRContentWithClaude(config, diff)
}

// GetProvider returns the appropriate provider based on configuration
// Supports "claude" (Claude Code CLI) and "openai" providers
func GetProvider(config *Config) (Provider, error) {
	providerName := "claude" // Default to Claude Code
	
	if providerValue, ok := config.MainConfig["llm_provider"].(string); ok && providerValue != "" {
		providerName = providerValue
	}

	switch providerName {
	case "openai":
		return &OpenAIProvider{}, nil
	case "claude":
		return &ClaudeProvider{}, nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s. Supported providers: openai, claude", providerName)
	}
}

// GeneratePRContentWithProvider generates PR content using the configured provider
func GeneratePRContentWithProvider(config *Config, diff string) (title, body string, err error) {
	provider, err := GetProvider(config)
	if err != nil {
		return "", "", err
	}

	return provider.GeneratePRContent(config, diff)
}