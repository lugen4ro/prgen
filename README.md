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

### Default Configuration Values

#### `config.json`

```json
{
  "llm_provider": "claude",
  "model": "claude-3-5-sonnet-20241022",
  "temperature": 0.7,
  "max_tokens": 2000
}
```

#### `title_instructions.md`

```markdown
# PR Title Generation Instructions

## Format Guidelines

- Follow conventional commit format: type(scope): description
- Keep titles under 72 characters
- Use imperative mood ("add" not "added")
- Be specific but concise

## Types

- feat: New feature
- fix: Bug fix
- docs: Documentation changes
- style: Code style changes (formatting, etc.)
- refactor: Code refactoring
- test: Adding or updating tests
- chore: Maintenance tasks

## Examples

- Good: "feat(auth): add JWT token validation middleware"
- Good: "fix(api): handle null responses in user endpoint"
- Bad: "update stuff"
- Bad: "fixed the bug"

## Scope Guidelines

- Use scope when the change affects a specific component
- Common scopes: api, ui, auth, db, config, docs
- Omit scope for global changes
```

#### `body_instructions.md`

```markdown
# PR Body Generation Instructions

## Style Guidelines

- Use clear, concise language
- Include relevant context and motivation
- Keep descriptions focused and actionable
- Use proper markdown formatting

## Content Requirements

- Always include a clear description of what changed
- List specific files or components affected
- Mention any breaking changes
- Include testing approach or validation steps
- Reference related issues or tickets when applicable

## Technical Notes

- Reference specific files when relevant
- Include code snippets if helpful for understanding
- Consider the target audience (reviewers)
- Keep the body informative but not overly verbose

## Examples

- Good description: "Refactored user authentication to use JWT tokens instead of sessions, improving scalability and security"
- Bad description: "changed auth stuff"
```

#### `title_example.md`

```
feat(auth): add JWT token validation middleware

fix(api): handle null responses in user endpoint

docs: update installation guide for Docker setup

refactor(db): optimize user query performance

test(auth): add integration tests for login flow
```

#### `body_example.md`

```markdown
## Description

Added JWT token validation middleware to improve API security and enable stateless authentication. This replaces the previous session-based approach and provides better scalability for our microservices architecture.

## Changes

- Implemented JWT validation middleware in auth/middleware.go
- Updated user authentication endpoints to issue JWT tokens
- Refactored protected routes to use token validation
- Added token refresh mechanism
- Removed session storage dependencies

## Files Modified

- `auth/middleware.go`: New JWT validation middleware
- `handlers/auth.go`: Updated login/logout to use JWT
- `routes/api.go`: Applied middleware to protected routes
- `config/auth.go`: JWT configuration settings

## Testing

- Added unit tests for JWT middleware (auth/middleware_test.go)
- Integration tests for login/logout flow
- Manual testing with Postman for token validation
- Load testing shows 40% improvement in response times

## Additional Notes

- Breaking change: clients must update to use Authorization header
- Migration guide added to docs/jwt-migration.md
- Tokens expire after 24 hours with 7-day refresh window
```

## Roadmap

- [ ] Add version command
- [ ] Edit generated PR before pushing to GitHub
- [ ] Fix preview layout shift for PRs with long text
