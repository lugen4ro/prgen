"""Base classes and interfaces for LLM providers."""

from abc import ABC, abstractmethod
from enum import Enum

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
