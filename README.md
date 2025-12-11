`prgen` is a tool that provides a flow for auto-generating GitHub PRs using Claude Code.

## Installation

Can be installed with brew.

```bash
brew install lugen4ro/homebrew-lugen4ro/prgen
```

### For MacOS User

Since the binary hasn't been notarized by Apple, macOS Gatekeeper may block it. You can bypass this by either:

1. Using the `--no-quarantine` flag during installation (will be deprecated in late 2026):

```bash
brew install --no-quarantine lugen4ro/homebrew-lugen4ro/prgen
```

2. Removing the quarantine attribute manually after installation like:

```bash
sudo xattr -rd com.apple.quarantine /opt/homebrew/bin/prgen
```

## Requirements

- Claude Code should be installed and logged in to.

## Usage

Just execute the `prgen` command in a checked out repository.

## Configuration

This tool uses config files placed under `~/.config/prgen/` which include:

### Core Configuration files

- `config.json` - Main configuration file with LLM settings
- `body_instructions.md` - PR body generation instructions for the LLM
- `title_instructions.md` - PR title generation instructions for the LLM

### Examples (for reference)

- `body_example.md` - Example of a well-formatted PR body
- `title_example.md` - Examples of good PR titles

All files are created automatically with sensible defaults on first run. Edit them to customize your PR generation style.
