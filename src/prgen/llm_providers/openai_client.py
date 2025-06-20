"""OpenAI provider implementation."""

import os

import instructor

from .base import BaseLLMClient, PRResponse


class OpenAIClient(BaseLLMClient):
    """OpenAI provider using Instructor."""

    def __init__(self):
        try:
            import openai

            self._openai = openai
        except ImportError:
            raise ImportError("openai package not installed. Install with: pip install openai") from None

    def is_available(self) -> bool:
        return bool(os.getenv("OPENAI_API_KEY"))

    def generate_pr_content(self, prompt: str) -> PRResponse:
        api_key = os.getenv("OPENAI_API_KEY")
        if not api_key:
            raise ValueError("OPENAI_API_KEY environment variable not set")

        client = instructor.from_openai(self._openai.OpenAI(api_key=api_key))

        response = client.chat.completions.create(
            model="gpt-4o-mini",
            response_model=PRResponse,
            messages=[{"role": "user", "content": prompt}],
            max_tokens=500,
            temperature=0.7,
        )
        return response
