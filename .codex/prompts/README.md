# Codex Custom Prompts

This directory contains custom prompts for [Codex](https://github.com/stackblitz/codex), converted from Claude Code custom slash commands.

## Overview

Codex is a CLI tool that provides AI-powered assistance for development tasks. Custom prompts allow you to create reusable, specialized commands that can be invoked with a simple slash command syntax (e.g., `/create-github-issue`).

## Directory Structure

```
~/.codex/prompts/
├── README.md                    # This file
├── create-github-issue.md       # Custom prompt for creating GitHub issues
└── [other-prompts].md           # Additional custom prompts
```

## Converting Claude Commands to Codex Prompts

### Automatic Conversion

Use the Claude Code conversion command to automatically convert Claude custom slash commands to Codex format:

```bash
/convert-custom-slash-for-codex <command-name>
```

**Example:**
```bash
/convert-custom-slash-for-codex create-github-issue
```

This command will:
1. Read the source command from `.claude/commands/`
2. Convert it to Codex format
3. Add proper YAML frontmatter
4. Handle parameter placeholders
5. Escape literal dollar signs
6. Write the output to `~/.codex/prompts/`

### Manual Conversion Process

If you prefer to convert manually or need to understand the process:

#### 1. YAML Frontmatter

Every Codex prompt must start with YAML frontmatter:

```yaml
---
description: Short description of the command (100-150 chars)
argument-hint: PARAM1=<value> PARAM2=<value>  # Optional, only if command takes parameters
---
```

**Examples:**

```yaml
---
description: Fix Issue $ISSUE_NUMBER
argument-hint: ISSUE_NUMBER=#123
---
```

```yaml
---
description: Create GitHub Issue
---
```

#### 2. Parameter Conversion

Convert Claude parameter placeholders to Codex format:

| Claude Format | Codex Format | Example |
|--------------|--------------|---------|
| `#{param_name}` | `$PARAM_NAME` | `#{issue_number}` → `$ISSUE_NUMBER` |
| `{param_name}` | `$PARAM_NAME` | `{file_path}` → `$FILE_PATH` |
| `{param-name}` | `$PARAM_NAME` | `{pr-number}` → `$PR_NUMBER` |

**Rules:**
- Convert to UPPERCASE
- Replace hyphens with underscores
- Preserve underscores in parameter names

**Example conversions:**

```markdown
# Claude format
Fix GitHub issue #{issue_number} by analyzing the code...

# Codex format
Fix GitHub issue $ISSUE_NUMBER by analyzing the code...
```

#### 3. Dollar Sign Escaping

Codex uses `$` for placeholders, so literal dollar signs must be escaped:

| Context | Original | Escaped |
|---------|----------|---------|
| Shell variable | `$HOME` | `$$HOME` |
| Shell command | `export PATH=$PATH:~/bin` | `export PATH=$$PATH:~/bin` |
| Price/currency | `$100` | `$$100` |
| Codex placeholder | `$ISSUE_NUMBER` | Keep as-is |

**Important:** Only escape literal dollar signs, NOT Codex placeholders!

#### 4. Preserve Markdown Structure

Maintain all markdown formatting:
- Headings (`#`, `##`, `###`)
- Lists (ordered and unordered)
- Code blocks with language tags
- Tables
- Links
- Emphasis (`*italic*`, `**bold**`)

## Prompt Format Specification

### Complete Example

```markdown
---
description: Fix Issue $ISSUE_NUMBER
argument-hint: ISSUE_NUMBER=#123
---
# Fix Issue $ISSUE_NUMBER

You are tasked with fixing GitHub issue $ISSUE_NUMBER for the go-crypto-wallet repository.

## Workflow

1. Fetch issue details:
   ```bash
   gh issue view $ISSUE_NUMBER
   ```

2. Analyze the issue and propose a fix

3. Implement the fix following project guidelines

4. Run verification commands:
   ```bash
   make go-lint
   make gotest
   ```

5. Commit and push changes

## Safety Rules

- Never include sensitive information (private keys, API tokens)
- Follow Clean Architecture principles
- Run tests before committing
```

### Field Descriptions

#### `description` (required)
- Short, clear description of what the command does
- Length: 100-150 characters recommended
- Displayed in Codex menu/help
- Can include parameter placeholders (e.g., `$ISSUE_NUMBER`)

#### `argument-hint` (optional)
- Only include if command accepts parameters
- Format: `PARAM1=<value> PARAM2=<value>`
- Use square brackets for optional parameters: `[PARAM=<value>]`
- Examples:
  - Single parameter: `ISSUE_NUMBER=#123`
  - Multiple parameters: `ISSUE_NUMBER=#123 PR_NUMBER=#456`
  - Optional parameter: `[FILE_PATH=<path>]`

## Available Prompts

### `create-github-issue`

Create well-structured GitHub issues for the go-crypto-wallet repository.

**Usage:**
```bash
codex /create-github-issue
```

**Workflow:**
1. Analyzes your request and creates a comprehensive issue proposal
2. Presents proposal for review
3. After approval, creates the issue using `gh` CLI

**No parameters required** - the command will gather information interactively.

## Creating New Prompts

### From Existing Claude Commands

1. **Check available Claude commands:**
   ```bash
   ls .claude/commands/
   ```

2. **Convert using the conversion command:**
   ```bash
   /convert-custom-slash-for-codex <command-name>
   ```

3. **Verify the conversion:**
   ```bash
   cat ~/.codex/prompts/<command-name>.md
   ```

4. **Test in Codex:**
   ```bash
   codex /<command-name>
   ```

### From Scratch

1. **Create a new markdown file:**
   ```bash
   touch ~/.codex/prompts/my-command.md
   ```

2. **Add YAML frontmatter:**
   ```yaml
   ---
   description: My custom command description
   argument-hint: PARAM=<value>  # If needed
   ---
   ```

3. **Write the prompt content:**
   - Clear instructions
   - Step-by-step workflow
   - Code examples
   - Safety rules
   - Error handling

4. **Test the prompt:**
   ```bash
   codex /my-command
   ```

## Best Practices

### 1. Clear Instructions

Write prompts that are:
- **Specific**: Clearly define what the command should do
- **Actionable**: Provide step-by-step workflows
- **Complete**: Include all necessary context and constraints

### 2. Parameter Design

- Use descriptive parameter names (`ISSUE_NUMBER`, not `NUM`)
- Provide clear examples in `argument-hint`
- Document required vs. optional parameters
- Use consistent naming conventions (UPPERCASE_WITH_UNDERSCORES)

### 3. Safety and Security

Always include safety rules:
- Never log/commit sensitive information
- Verify operations before execution
- Follow project-specific guidelines
- Include validation steps

### 4. Project Context

For project-specific commands:
- Reference relevant documentation (CLAUDE.md, agents/*.md, etc.)
- Include architecture guidelines
- Specify required tools and versions
- Link to related resources

### 5. Error Handling

Include guidance for common errors:
- Missing tools or dependencies
- Authentication failures
- Invalid parameters
- Network issues

## Testing Prompts

### 1. Syntax Validation

Check YAML frontmatter is valid:
```bash
# Should not show YAML errors
head -10 ~/.codex/prompts/<command>.md
```

### 2. List Prompts

Verify prompt appears in Codex:
```bash
codex list
```

### 3. Dry Run

Test the prompt without making changes:
```bash
codex /<command>
```

Review the AI's response to ensure:
- Prompt loads correctly
- Parameters are substituted properly
- Instructions are clear and followed
- Output format is as expected

### 4. Full Test

Run the prompt for a real task:
- Verify all steps execute correctly
- Check that safety rules are followed
- Confirm output matches expectations

## Troubleshooting

### Prompt Not Appearing in List

**Problem:** `codex list` doesn't show your custom prompt

**Solutions:**
1. Check file is in `~/.codex/prompts/` directory
2. Verify filename ends with `.md`
3. Check YAML frontmatter is valid
4. Restart Codex if running in interactive mode

### YAML Parsing Errors

**Problem:** Errors when loading prompt

**Solutions:**
1. Verify YAML frontmatter starts and ends with `---`
2. Check for proper indentation (use spaces, not tabs)
3. Ensure `description:` field is present
4. Remove any special characters in YAML values

### Parameter Not Substituted

**Problem:** `$PARAM` appears literally instead of being substituted

**Solutions:**
1. Check parameter is specified in `argument-hint:`
2. Verify parameter name matches exactly (case-sensitive)
3. Ensure you're passing the parameter when invoking the command
4. Check for typos in parameter names

### Dollar Signs Displaying Incorrectly

**Problem:** Shell variables or prices showing as `$$VAR` or `$$100`

**Solutions:**
- This is correct! `$$` is how you write literal `$` in Codex prompts
- Codex will display them correctly as single `$` when executing
- Only Codex placeholders (`$PARAM_NAME`) use single `$`

## Maintenance

### Updating Prompts

When updating Claude commands:

1. **Update the source:** Edit `.claude/commands/<command>.md`
2. **Reconvert:** Run `/convert-custom-slash-for-codex <command>`
3. **Review changes:** Check the updated `~/.codex/prompts/<command>.md`
4. **Test:** Verify the updated prompt works correctly

### Version Control

Consider tracking prompts in version control:

```bash
# Add prompts directory to your dotfiles repo
cd ~/.codex
git init
git add prompts/
git commit -m "Add Codex custom prompts"
```

### Syncing Across Machines

Share prompts between development environments:

1. **Commit to dotfiles repo:**
   ```bash
   cd ~/.codex
   git add prompts/
   git commit -m "Update custom prompts"
   git push
   ```

2. **Pull on other machines:**
   ```bash
   cd ~/.codex
   git pull
   ```

## Examples

### Parameterless Command

```markdown
---
description: Run linting and fix issues
---
# Fix Linting Issues

Run golangci-lint and automatically fix issues in the codebase.

## Steps

1. Run linter: `make go-lint`
2. Review and fix issues
3. Verify fixes: `make go-lint`
```

### Command with Single Parameter

```markdown
---
description: Fix Issue $ISSUE_NUMBER
argument-hint: ISSUE_NUMBER=#123
---
# Fix Issue $ISSUE_NUMBER

Fix GitHub issue $ISSUE_NUMBER for this repository.

## Usage

When you run this command, provide the issue number:
```bash
codex /fix-issue ISSUE_NUMBER=#456
```
```

### Command with Multiple Parameters

```markdown
---
description: Review PR $PR_NUMBER and address comments
argument-hint: PR_NUMBER=#123 [BRANCH=<name>]
---
# Fix PR Review Comments

Address review comments for PR $PR_NUMBER.

## Parameters

- `PR_NUMBER` (required): Pull request number (e.g., #123)
- `BRANCH` (optional): Branch name (defaults to current branch)

## Steps

1. Fetch PR details: `gh pr view $PR_NUMBER`
2. Review comments and suggestions
3. Implement fixes
4. Push changes
```

## Resources

### Documentation

- **Codex GitHub**: https://github.com/stackblitz/codex
- **Claude Code**: https://claude.com/claude-code
- **Project Guidelines**: See `CLAUDE.md` in the project root

### Related Files

- **Claude commands**: `.claude/commands/` (source)
- **Codex prompts**: `~/.codex/prompts/` (converted)
- **Project instructions**: `CLAUDE.md`, `agents/*.md`

### Conversion Command

Location: `.claude/commands/convert-custom-slash-for-codex.md`

Run in Claude Code:
```bash
/convert-custom-slash-for-codex <command-name>
```

## Contributing

When creating or updating prompts:

1. Follow the format specification
2. Include comprehensive documentation
3. Add safety rules and error handling
4. Test thoroughly before committing
5. Update this README if adding new conventions

## License

These prompts are specific to the go-crypto-wallet project and follow the same license as the main project.
