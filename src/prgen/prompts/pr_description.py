PROMPT_TEMPLATE = """
Generate a PR title and description based on this git diff.

For the title:
- Keep it 50-72 characters long
- Start with a conventional commit type (feat, fix, refactor, docs, style, test, chore)
- Use imperative mood (e.g., "add", "fix", "update")
- Clearly summarize the main change

For the description:
- Provide a brief summary of changes
- List key improvements or fixes
- Mention any breaking changes if applicable

Git diff:
{diff}
"""


def get_pr_content_prompt(diff: str) -> str:
    """
    Generate PR title and description prompt with the provided git diff.

    Args:
        diff (str): The git diff content to analyze

    Returns:
        str: The formatted prompt ready for LLM consumption
    """
    return PROMPT_TEMPLATE.format(diff=diff)
