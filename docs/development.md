# Development Guide

## Prerequisites

- Go 1.21+
- [Task](https://taskfile.dev/) - Task runner (`brew install go-task`)
- [GoReleaser](https://goreleaser.com/) - Release automation (`brew install goreleaser`)

## Available Tasks

Run `task --list` to see all available commands:

```bash
task build          # Build the Go binary
task install        # Build and install globally to /usr/local/bin
task clean          # Remove build artifacts
task release        # Create a new release
task release:dry-run # Test release build without publishing
task release:check  # Validate GoReleaser configuration
```

## Building

```bash
# Build locally
task build

# Install to /usr/local/bin
task install
```

## Releasing a New Version

### Prerequisites

1. Ensure `GITHUB_TOKEN` is set with repo permissions:
   ```bash
   export GITHUB_TOKEN=ghp_xxxxx
   ```

2. Ensure all changes are committed and pushed to main

### Release Process

```bash
# 1. (Optional) Validate GoReleaser config
task release:check

# 2. (Optional) Test the release build
task release:dry-run

# 3. Create and publish the release
task release VERSION=0.2.0
```

This will:
1. Create a git tag `v0.2.0`
2. Push the tag to origin
3. Build binaries for darwin/linux (amd64/arm64)
4. Create a universal binary for macOS
5. Upload release artifacts to GitHub
6. Auto-update the Homebrew cask at [lugen4ro/homebrew-lugen4ro](https://github.com/lugen4ro/homebrew-lugen4ro)

### Installing via Homebrew

Users can install with:
```bash
brew tap lugen4ro/lugen4ro
brew install --cask prgen
```
