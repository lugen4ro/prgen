#!/bin/bash
# Local CI script that runs the same checks as GitHub Actions

set -e  # Exit on any error

echo "🔧 Running local CI checks (same as GitHub Actions)"
echo "================================================="

echo ""
echo "📦 Installing dependencies..."
python -m pip install --upgrade pip
pip install -e ".[dev]"

echo ""
echo "🧹 Running ruff linting..."
ruff check .

echo ""
echo "🎨 Checking code formatting with ruff..."
ruff format --check .

echo ""
echo "🔍 Running PyRight type checking..."
pyright

echo ""
echo "🧪 Running tests with pytest..."
pytest

echo ""
echo "✅ All CI checks passed! Ready to commit and push."