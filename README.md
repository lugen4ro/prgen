
`prgen` is a tool for auto-generating GitHub PRs with the power of AI.

## Installation

Can be installed with brew.
Requires the `--no-quarantine` flag for macOS because it hasn't been notarized by Apple.

```bash
brew install --no-quarantine lugen4ro/homebrew-lugen4ro/prgen
```

## Usage

Just execute the `prgen` command in your git repository.

## Configuration

This tool uses config files placed under `~/.config/prgen/` which include:

### Core Configuration
- `config.json` - Main configuration file with LLM settings
- `body_instructions.md` - PR body generation instructions for the LLM
- `title_instructions.md` - PR title generation instructions for the LLM

#### AI Provider Configuration

The tool supports two AI providers:

**Claude Code CLI (Default)**
- Set `"llm_provider": "claude"` in config.json
- Requires Claude Code CLI to be installed and available in PATH
- Uses the local Claude Code CLI for generation

**OpenAI**
- Set `"llm_provider": "openai"` in config.json  
- Requires `OPENAI_API_KEY` environment variable
- Supports models like "gpt-4", "gpt-3.5-turbo"

### Examples (for reference)
- `body_example.md` - Example of a well-formatted PR body
- `title_example.md` - Examples of good PR titles

All files are created automatically with sensible defaults on first run. Edit them to customize your PR generation style.
