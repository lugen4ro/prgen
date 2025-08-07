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