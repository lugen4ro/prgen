"""Factory for creating LLM clients."""

from typing import Optional

from .anthropic_client import AnthropicClient
from .base import BaseLLMClient, LLMProvider
from .google_client import GoogleClient
from .openai_client import OpenAIClient


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
