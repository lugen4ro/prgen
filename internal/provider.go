package internal

// RefinementContext holds information needed for refining a previously generated PR
type RefinementContext struct {
	SessionID string // Session ID for continuing the conversation
	Feedback  string // User's feedback for refinement
}

// Provider represents an AI provider interface
type Provider interface {
	GeneratePRContent(config *Config, diff, background string) (*PRGenerationResult, error)
	RefinePRContent(config *Config, diff, background string, refinement *RefinementContext) (*PRGenerationResult, error)
}

// ClaudeProvider implements the Provider interface for Claude Code CLI
type ClaudeProvider struct{}

// GeneratePRContent generates PR content using Claude Code CLI
func (p *ClaudeProvider) GeneratePRContent(config *Config, diff, background string) (*PRGenerationResult, error) {
	return GeneratePRContentWithClaude(config, diff, background, nil)
}

// RefinePRContent refines PR content based on user feedback using Claude Code CLI
func (p *ClaudeProvider) RefinePRContent(config *Config, diff, background string, refinement *RefinementContext) (*PRGenerationResult, error) {
	return GeneratePRContentWithClaude(config, diff, background, refinement)
}

// GetProvider returns the Claude provider
func GetProvider(config *Config) (Provider, error) {
	return &ClaudeProvider{}, nil
}

// GeneratePRContentWithProvider generates PR content using the Claude provider
func GeneratePRContentWithProvider(config *Config, diff, background string) (*PRGenerationResult, error) {
	provider, err := GetProvider(config)
	if err != nil {
		return nil, err
	}

	return provider.GeneratePRContent(config, diff, background)
}

// RefinePRContentWithProvider refines PR content using the Claude provider
func RefinePRContentWithProvider(config *Config, diff, background string, refinement *RefinementContext) (*PRGenerationResult, error) {
	provider, err := GetProvider(config)
	if err != nil {
		return nil, err
	}

	return provider.RefinePRContent(config, diff, background, refinement)
}
