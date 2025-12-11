package internal

// Provider represents an AI provider interface
type Provider interface {
	GeneratePRContent(config *Config, diff, background string) (title, body string, err error)
}

// ClaudeProvider implements the Provider interface for Claude Code CLI
type ClaudeProvider struct{}

// GeneratePRContent generates PR content using Claude Code CLI
func (p *ClaudeProvider) GeneratePRContent(config *Config, diff, background string) (title, body string, err error) {
	return GeneratePRContentWithClaude(config, diff, background)
}

// GetProvider returns the Claude provider
func GetProvider(config *Config) (Provider, error) {
	return &ClaudeProvider{}, nil
}

// GeneratePRContentWithProvider generates PR content using the Claude provider
func GeneratePRContentWithProvider(config *Config, diff, background string) (title, body string, err error) {
	provider, err := GetProvider(config)
	if err != nil {
		return "", "", err
	}

	return provider.GeneratePRContent(config, diff, background)
}
