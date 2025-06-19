"""Google provider implementation."""

import asyncio
import os

import instructor

from .base import BaseLLMClient, PRResponse


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
