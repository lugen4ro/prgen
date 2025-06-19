import subprocess


def get_diff() -> str | None:
    """
    Get the git diff between the current branch and the default target branch.

    Returns:
        Optional[str]: The git diff output if changes exist, None otherwise.
                      Also prints status messages about the comparison.

    Raises:
        subprocess.CalledProcessError: If any git command fails.
    """

    current_branch: str = subprocess.run(
        ["git", "branch", "--show-current"], capture_output=True, text=True, check=True
    ).stdout.strip()

    target_branch: str = subprocess.run(
        ["git", "config", "--get", "init.defaultBranch"],
        capture_output=True,
        text=True,
        check=True,
    ).stdout.strip()

    diff: str = subprocess.run(
        ["git", "diff", f"origin/{target_branch}...{current_branch}"],
        capture_output=True,
        text=True,
        check=True,
    ).stdout.strip()

    # Filter out only the actual code changes (lines starting with + or -)
    # but exclude the diff header lines that also start with --- or +++
    code_changes = [
        line
        for line in diff.splitlines()
        if (line.startswith("+") or line.startswith("-")) and not (line.startswith("+++") or line.startswith("---"))
    ]

    if diff:
        print(f"Getting code changes between '{current_branch}' and '{target_branch}' ({len(code_changes)} lines)")
        return diff
    else:
        print(f"Empty diff  between '{current_branch}' and '{target_branch}'.")
        return None
