from .llm_providers import LLMClientFactory, LLMProvider
from .prompts.pr_description import get_pr_description_prompt


def generate_pr_description(diff: str, provider: LLMProvider | None = None) -> str:
    """
    Generate PR description using specified or auto-detected LLM provider.

    Args:
        diff: Git diff content to analyze
        provider: Optional specific provider to use. If None, auto-detect.

    Returns:
        Generated PR description string
    """
    prompt = get_pr_description_prompt(diff)

    client = LLMClientFactory.create_client(provider)
    response = client.generate_pr_description(prompt)

    # Format the response as a single string for backward compatibility
    return f"{response.title}\n\n{response.body}"
