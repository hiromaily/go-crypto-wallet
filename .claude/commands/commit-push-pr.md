# Commit, Push, and Create Pull Request

You are tasked with creating a commit, pushing changes, and creating a pull request based on the modifications made during this conversation.

## Context Analysis

First, determine the source of change information:

1. **Check git status**: Run `git status` and `git diff` to see actual file changes
2. **Review conversation context**: Consider what was discussed and modified in this session
3. **Combine both sources**: Use git changes as the primary source, but leverage conversation context for understanding the "why" behind changes

## Commit Message Guidelines

Follow conventional commit format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type Options

- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring
- `style`: Code style changes (formatting, lint fixes)
- `docs`: Documentation changes
- `test`: Test additions or modifications
- `chore`: Maintenance tasks
- `perf`: Performance improvements

### Commit Message Best Practices

- **Subject line**: 50 characters or less, imperative mood ("fix" not "fixed")
- **Body**: Explain what and why, not how (wrap at 72 characters)
- **Be specific**: Reference issue numbers, explain context
- **Atomic commits**: If multiple unrelated changes exist, create separate commits

## Pull Request Guidelines

### PR Title

- Should match the commit message subject (or summarize multiple commits)
- Clear and descriptive

### PR Description Structure

```markdown
## Summary
Brief overview of what this PR accomplishes

## Changes
- List of key changes made
- Group related changes together

## Motivation
Why these changes were necessary

## Testing
- How the changes were tested
- Any manual testing steps performed

## Related Issues
Fixes #123, Relates to #456
```

## Prerequisites Check

Before starting, verify all required tools are installed:

1. **Git**: Check with `git --version`
   - Required for all version control operations
   - If missing: Install from <https://git-scm.com/>

2. **GitHub CLI (gh)**: Check with `gh --version`
   - Required for creating pull requests via CLI
   - If missing: Install from <https://cli.github.com/>
   - Alternative: Provide manual PR creation instructions

3. **Git authentication**: Verify with `git remote -v` and `gh auth status`
   - Ensure remote repository is accessible
   - Ensure GitHub CLI is authenticated

**If any required tool is missing, stop and display an error message with installation instructions. Do not proceed with the workflow.**

## Process

1. **Check current branch**:
   - Run `git branch --show-current` to identify the current branch
   - If on `main` or `master`, proceed to step 2
   - If already on a feature branch, skip to step 3

2. **Create feature branch** (if on main/master):
   - Analyze the changes to determine appropriate branch name
   - Branch naming convention: `feature/<descriptive-name>`
     - `feature/fix-lint-errors` for lint fixes
     - `feature/add-user-authentication` for new features
     - `feature/update-dependencies` for dependency updates
     - `feature/refactor-api-layer` for refactoring
   - Use kebab-case (lowercase with hyphens)
   - Keep it concise but descriptive (2-4 words recommended)
   - Create and switch to the new branch: `git checkout -b feature/<name>`

3. **Analyze changes**:
   - Run `git status` and `git diff` to review all modifications
   - Identify the scope and type of changes

4. **Stage and commit**:
   - Stage relevant files with `git add`
   - Create a well-formatted commit message
   - If changes are large and diverse, consider multiple atomic commits

5. **Push changes**:
   - Push to remote repository
   - Set upstream if this is a new branch

6. **Create Pull Request**:
   - Use GitHub CLI (`gh pr create`) or provide instructions for manual creation
   - Include comprehensive PR description
   - Add appropriate labels if applicable
   - Request reviewers if known

## Output Format

Before executing, show the plan:

1. **Proposed branch name** (if creating new branch)
2. **Proposed commit message(s)**
3. **Proposed PR title and description**
4. **List of files to be committed**

Ask for confirmation before proceeding with the actual git operations.

## Safety Checks

- Verify we're not on a protected branch (main/master) before committing
- Confirm no sensitive information is being committed
- Ensure all files are intended for this commit
- Check that tests pass (if applicable) before creating PR

## Error Handling

If any step fails:

- Explain what went wrong clearly
- Provide remediation steps
- Don't proceed to subsequent steps until the issue is resolved
