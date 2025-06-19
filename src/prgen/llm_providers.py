"""
LLM provider implementations using Instructor for structured outputs.
"""

import asyncio
import os
from abc import ABC, abstractmethod
from enum import Enum
from typing import Optional

import instructor
from pydantic import BaseModel


class LLMProvider(str, Enum):
    """Available LLM providers."""

    OPENAI = "openai"
    ANTHROPIC = "anthropic"
    GOOGLE = "google"


class PRResponse(BaseModel):
    """Structured response for PR generation."""

    title: str
    body: str


class BaseLLMClient(ABC):
    """Base class for LLM clients."""

    @abstractmethod
    def generate_pr_description(self, prompt: str) -> PRResponse:
        """Generate PR description from prompt."""
        pass

    @abstractmethod
    def is_available(self) -> bool:
        """Check if the provider is available (API key exists)."""
        pass


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

    def generate_pr_description(self, prompt: str) -> PRResponse:
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

    def generate_pr_description(self, prompt: str) -> PRResponse:
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


class GoogleClient(BaseLLMClient):
    """Google (Gemini) provider using Instructor."""

    def __init__(self):
        try:
            import google.generativeai as genai

            self._genai = genai
        except ImportError:
            raise ImportError(
                "google-generativeai package not installed. Install with: pip install google-generativeai"
            ) from None

    def is_available(self) -> bool:
        return bool(os.getenv("GOOGLE_API_KEY"))

    def generate_pr_description(self, prompt: str) -> PRResponse:
        api_key = os.getenv("GOOGLE_API_KEY")
        if not api_key:
            raise ValueError("GOOGLE_API_KEY environment variable not set")

        # Import at runtime to avoid pyright export issues
        try:
            # Import from the correct modules as suggested by pyright
            from google.generativeai.client import configure
            from google.generativeai.generative_models import GenerativeModel
        except ImportError:
            raise ImportError("google-generativeai package not installed") from None

        configure(api_key=api_key)
        client = instructor.from_gemini(
            client=GenerativeModel("gemini-1.5-flash"),
            mode=instructor.Mode.GEMINI_JSON,
        )

        # The instructor client may return a coroutine, so we need to handle it
        result = client.chat.completions.create(
            response_model=PRResponse,
            messages=[{"role": "user", "content": prompt}],
        )

        # If it's a coroutine, run it synchronously
        if asyncio.iscoroutine(result):
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
            try:
                response: PRResponse = loop.run_until_complete(result)
            finally:
                loop.close()
        else:
            response = result

        return response


class LLMClientFactory:
    """Factory for creating LLM clients."""

    _providers = {
        LLMProvider.OPENAI: OpenAIClient,
        LLMProvider.ANTHROPIC: AnthropicClient,
        LLMProvider.GOOGLE: GoogleClient,
    }

    @classmethod
    def create_client(cls, provider: Optional[LLMProvider] = None) -> BaseLLMClient:
        """
        Create an LLM client.

        Args:
            provider: Specific provider to use. If None, auto-detect based on available API keys.

        Returns:
            Configured LLM client.

        Raises:
            ValueError: If no provider is available or specified provider is unavailable.
        """
        if provider:
            client_class = cls._providers.get(provider)
            if not client_class:
                raise ValueError(f"Unknown provider: {provider}")

            client = client_class()
            if not client.is_available():
                raise ValueError(f"Provider {provider} is not available (missing API key)")
            print(f"Using LLM provider: {provider.value.title()}")
            return client

        # Auto-detect provider
        for provider_enum, client_class in cls._providers.items():
            try:
                client = client_class()
                if client.is_available():
                    print(f"Using LLM provider: {provider_enum.value.title()}")
                    return client
            except ImportError:
                continue

        raise ValueError(
            "No LLM provider available. Set one of:\n"
            "• OPENAI_API_KEY (get from: https://platform.openai.com/api-keys)\n"
            "• ANTHROPIC_API_KEY (get from: https://console.anthropic.com/)\n"
            "• GOOGLE_API_KEY (get from: https://aistudio.google.com/app/apikey)"
        )

    @classmethod
    def get_available_providers(cls) -> list[LLMProvider]:
        """Get list of available providers based on installed packages and API keys."""
        available = []
        for provider_enum, client_class in cls._providers.items():
            try:
                client = client_class()
                if client.is_available():
                    available.append(provider_enum)
            except ImportError:
                continue
        return available
