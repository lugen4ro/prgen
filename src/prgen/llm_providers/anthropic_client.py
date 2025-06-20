"""Anthropic provider implementation."""

import os

import instructor

from .base import BaseLLMClient, PRResponse


class AnthropicClient(BaseLLMClient):
    """Anthropic (Claude) provider using Instructor."""

    def __init__(self):
        try:
            import anthropic

            self._anthropic = anthropic
        except ImportError:
            raise ImportError("anthropic package not installed. Install with: pip install anthropic") from None

    def is_available(self) -> bool:
        return bool(os.getenv("ANTHROPIC_API_KEY"))

    def generate_pr_content(self, prompt: str) -> PRResponse:
        api_key = os.getenv("ANTHROPIC_API_KEY")
        if not api_key:
            raise ValueError("ANTHROPIC_API_KEY environment variable not set")

        client = instructor.from_anthropic(self._anthropic.Anthropic(api_key=api_key))

        response = client.messages.create(
            model="claude-3-haiku-20240307",
            response_model=PRResponse,
            max_tokens=500,
            messages=[{"role": "user", "content": prompt}],
        )
        return response
