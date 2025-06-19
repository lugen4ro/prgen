# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## IMPORTANT

Adhere to ruff rules. Execute `ruff check .` before finishing and fix any issues.
- Execute pyright with ruff before finishing

## Project Overview

prgen is a CLI tool for generating GitHub PR titles & bodies using LLMs (OpenAI GPT or Google Gemini). The tool analyzes git diffs between the current branch and the default target branch to automatically generate PR descriptions.

**Status**: Under development, not yet in usable state.

## Commands

### Development & Quality Checks

```bash
# Install development dependencies
pip install -e ".[dev]"

# Linting and formatting
hatch run lint         # Check code with ruff
hatch run lint-fix     # Auto-fix linting issues
hatch run format       # Format code with ruff
hatch run typecheck    # Run pyright type checking
hatch run checks       # Run both typecheck and lint

# Manual tool usage
ruff check .           # Direct ruff linting
ruff check --fix .     # Direct ruff auto-fix
ruff format .          # Direct ruff formatting
pyright               # Direct type checking
```

### Installation & Usage

```bash
# Install from PyPI
pip install prgen

# Run the tool (must be in git repository on feature branch)
prgen
```

## Architecture

### Core Flow

1. **Diff extraction** (`src/prgen/diff.py`): Get git diff between current branch and default target branch
2. **LLM generation** (`src/prgen/llm_client.py`): Auto-detect API key and call appropriate LLM provider
3. **Prompt templating** (`src/prgen/prompts/pr_description.py`): Format diff data for LLM consumption

### Key Components

- **Entry point**: `src/prgen/main.py:main()` - orchestrates the full flow
- **Git integration**: Uses subprocess calls to git for branch detection and diff generation
- **Multi-provider LLM support**: Auto-detects OPENAI_API_KEY or GOOGLE_API_KEY environment variables
- **Template system**: Modular prompt templates in `src/prgen/prompts/`

### Dependencies

- **LLM clients**: openai, google-generativeai
- **Data validation**: pydantic
- **Code interaction**: instructor
- **Dev tools**: ruff (linting/formatting), pyright (type checking), hatch (project management)

### Configuration

- Line length: 120 characters
- Python version: >=3.9
- Type checking: basic mode with pyright
- Formatting: double quotes, space indentation
