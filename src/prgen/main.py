from rich.console import Console
from rich.markdown import Markdown
from rich.panel import Panel

from .diff import get_diff
from .llm_client import generate_pr_content


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
        title, description = generate_pr_content(diff)
    except Exception as e:
        print(f"Error: {e}")
        return

    # ╭─────────────────────────────────────────────────────────────────────
    # │ 3. Present generated content
    # ╰─────────────────────────────────────────────────────────────────────
    console = Console()

    # Display title
    title_panel = Panel(
        title,
        title="[bold cyan]PR Title[/bold cyan]",
        border_style="cyan",
        padding=(0, 1)
    )

    # Render description as markdown
    description_md = Markdown(description)
    description_panel = Panel(
        description_md,
        title="[bold cyan]PR Description[/bold cyan]",
        border_style="cyan",
        padding=(0, 1)
    )

    console.print("\n🚀 [bold green]Generated PR Content[/bold green]\n")
    console.print(title_panel)
    console.print()
    console.print(description_panel)

    # ╭─────────────────────────────────────────────────────────────────────
    # │ 4. Create PR on GitHub
    # ╰─────────────────────────────────────────────────────────────────────
