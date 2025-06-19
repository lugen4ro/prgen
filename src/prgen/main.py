from .diff import get_diff
from .llm_client import generate_pr_description


def main():
    # ╭─────────────────────────────────────────────────────────────────────
    # │ 1. Get diff
    # ╰─────────────────────────────────────────────────────────────────────
    diff = get_diff()

    if not diff:
        print("No changes to generate PR for.")
        return

    # ╭─────────────────────────────────────────────────────────────────────
    # │ 2. Generate PR content
    # ╰─────────────────────────────────────────────────────────────────────
    try:
        pr_description = generate_pr_description(diff)
        print(f"Generated PR Description:\n{pr_description}")
    except Exception as e:
        print(f"Error: {e}")
        return

    # ╭─────────────────────────────────────────────────────────────────────
    # │ 3. Propose PR content
    # ╰─────────────────────────────────────────────────────────────────────

    # ╭─────────────────────────────────────────────────────────────────────
    # │ 4. Create PR on GitHub
    # ╰─────────────────────────────────────────────────────────────────────
