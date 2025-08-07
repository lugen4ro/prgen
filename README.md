> [!CAUTION]
> This project is under development and as such is not yet in a usable state.

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
- `body_template.md` - PR body structure with instructional placeholders
- `title_template.md` - PR title format guidelines
- `body_instructions.md` - PR body generation instructions for the LLM
- `title_instructions.md` - PR title generation instructions for the LLM

### Examples (for reference)
- `body_example.md` - Example of a well-formatted PR body
- `title_example.md` - Examples of good PR titles

All files are created automatically with sensible defaults on first run. Edit them to customize your PR generation style.
