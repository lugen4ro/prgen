from .llm_providers import LLMClientFactory, LLMProvider
from .prompts.pr_description import get_pr_content_prompt


def generate_pr_content(diff: str, provider: LLMProvider | None = None) -> tuple[str, str]:
    """
    Generate PR title and description using specified or auto-detected LLM provider.

    Args:
        diff: Git diff content to analyze
        provider: Optional specific provider to use. If None, auto-detect.

    Returns:
        Tuple of (title, description)
    """
    prompt = get_pr_content_prompt(diff)

    client = LLMClientFactory.create_client(provider)
    response = client.generate_pr_content(prompt)

    return response.title, response.body
