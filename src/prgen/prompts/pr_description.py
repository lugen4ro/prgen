PROMPT_TEMPLATE = """
Generate a concise PR description based on this git diff. Include:
- A brief summary of changes
- Key improvements or fixes
- Any breaking changes if applicable

Git diff:
{diff}
"""


def get_pr_description_prompt(diff: str) -> str:
    """
    Generate PR description prompt with the provided git diff.

    Args:
        diff (str): The git diff content to analyze

    Returns:
        str: The formatted prompt ready for LLM consumption
    """
    return PROMPT_TEMPLATE.format(diff=diff)
