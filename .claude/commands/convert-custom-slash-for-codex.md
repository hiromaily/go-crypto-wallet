# Convert Claude Command to Codex Prompt

You are tasked with converting a Claude custom slash command from `.claude/commands/` to Codex-optimized prompt format for `~/.codex/prompts/`.

## Parameter

- **COMMAND_NAME** (required): Name of the Claude command to convert (without `.md` extension)
  - Example: `COMMAND_NAME=fix-issue`
  - Example: `COMMAND_NAME=fix-linter`

## Context

This repository uses both Claude (via `.claude/commands/`) and Codex (via `~/.codex/prompts/`). Codex prompts require:

1. **YAML frontmatter** with `description:` and `argument-hint:` fields
2. **Placeholder syntax** (`$1`, `$2`, `$ARGUMENTS`, `$NAMED_PARAM`) instead of markdown parameter references
3. **Literal dollar signs** using `$$` syntax

## Conversion Process

### Step 1: Validate Input

1. Check that `COMMAND_NAME` parameter is provided
2. Verify source file exists: `.claude/commands/$COMMAND_NAME.md`
3. If file doesn't exist, report error and list available commands

### Step 2: Read Source Command

1. Read the entire content from `.claude/commands/$COMMAND_NAME.md`
2. Store the original content for conversion

### Step 3: Extract Metadata

#### Description
- Extract from the first heading (without `#` symbols)
- If heading contains parameter placeholders, keep them for now (will convert later)
- Limit to 100-150 characters for Codex menu display
- Example: `# Fix Issue #{issue_number}` → Description: `Fix Issue #{issue_number}`

#### Argument Hints
Analyze the command content for parameter patterns:

**Common Parameter Patterns:**
- `#{issue_number}` or `{issue_number}` → `ISSUE_NUMBER=#123`
- `#{pr_number}` or `{pr_number}` → `PR_NUMBER=#456`
- `{file_path}` or `FILE=` → `FILE=<path>`
- `{name}` or `{description}` → `NAME=<value>` or `DESCRIPTION=<value>`
- `COMMAND_NAME=` → `COMMAND_NAME=<name>`

**Argument Hint Format:**
- Single parameter: `ISSUE_NUMBER=#123`
- Multiple parameters: `ISSUE_NUMBER=#123 PR_NUMBER=#456`
- Optional parameters: `[PARAM=<value>]` (use square brackets)
- Required parameters: `PARAM=<value>` (no brackets)

**Detection Rules:**
1. Look for patterns like `#{param}`, `{param}`, or `PARAM=` in the content
2. Convert parameter names to uppercase with underscores (e.g., `issue_number` → `ISSUE_NUMBER`)
3. Infer whether parameters are required or optional based on context:
   - If command says "required" or shows usage without brackets → required
   - If command says "optional" or shows `[param]` in usage → optional

### Step 4: Transform Content

#### 4.1 Generate YAML Frontmatter

```yaml
---
description: [extracted description]
argument-hint: [PARAM1=<value>] [PARAM2=<value>]
---
```

**Important:**
- If no parameters detected, omit `argument-hint:` field entirely
- Parameter order in `argument-hint:` should match the order they appear in the command

#### 4.2 Convert Parameter References

Replace parameter placeholders throughout the content:

**Conversion Rules:**
- `#{param_name}` → `$PARAM_NAME` (uppercase, underscores preserved)
- `{param_name}` → `$PARAM_NAME` (uppercase, underscores preserved)
- `{param-name}` → `$PARAM_NAME` (hyphens converted to underscores)
- Keep context: if parameter is in a sentence, maintain readability

**Common Conversions:**
- `#{issue_number}` → `$ISSUE_NUMBER`
- `#{pr_number}` → `$PR_NUMBER`
- `{file_path}` → `$FILE`
- `{name}` → `$NAME`
- `{description}` → `$DESCRIPTION`
- `{brief-description}` → `$BRIEF_DESCRIPTION`
- `COMMAND_NAME=` → remains as is (this is the argument hint syntax)

**Special Cases:**
- `$ARGUMENTS` placeholder: Keep as-is (Codex built-in)
- URLs with dollar signs: Convert `$` → `$$` except in placeholders
- Code examples showing `$` in shell commands: Convert to `$$`
- Placeholders already in correct format (`$PARAM`): Keep as-is

#### 4.3 Handle Dollar Signs

**Rules:**
1. **Preserve Codex placeholders**: `$PARAM_NAME`, `$1`, `$2`, `$ARGUMENTS` → Keep as-is
2. **Escape literal dollars**: Shell variables, prices, etc. → Convert `$` to `$$`
3. **In code blocks**: If showing shell examples, escape: `$HOME` → `$$HOME`
4. **In explanatory text**: Context-dependent - if referring to a placeholder, keep; if referring to shell variable, escape

**Examples:**
- Shell command example: `export PATH=$PATH:~/bin` → `export PATH=$$PATH:~/bin`
- Price mention: `$100` → `$$100`
- Codex placeholder: `$ISSUE_NUMBER` → Keep as-is
- Shell variable in text: "Set $HOME variable" → "Set $$HOME variable"

### Step 5: Preserve Markdown Structure

**Keep Intact:**
- All headings (`#`, `##`, `###`, etc.)
- Lists (numbered and bulleted)
- Code blocks with language tags
- Blockquotes
- Tables
- Links (preserve markdown link syntax)
- Emphasis (`*italic*`, `**bold**`)

**Maintain:**
- Proper indentation in lists and code blocks
- Line breaks and paragraph structure
- Code block language specifiers

### Step 6: Generate Output File

1. **Create directory if needed**:
   ```bash
   mkdir -p ~/.codex/prompts
   ```

2. **Construct output content**:
   ```markdown
   ---
   description: [extracted description]
   argument-hint: [PARAM=<value>]
   ---

   [converted command content]
   ```

3. **Write to file**:
   - Path: `~/.codex/prompts/$COMMAND_NAME.md`
   - Preserve line endings (LF on Unix/Mac, CRLF on Windows)
   - Use UTF-8 encoding

4. **Set permissions**:
   - Make file readable: `chmod 644 ~/.codex/prompts/$COMMAND_NAME.md`

### Step 7: Provide Output Summary

After successful conversion, display:

```
✅ Conversion completed successfully

Source: .claude/commands/$COMMAND_NAME.md
Output: ~/.codex/prompts/$COMMAND_NAME.md

Extracted metadata:
- Description: [description text]
- Parameters: [list of detected parameters]

Conversions applied:
- [list of parameter conversions, e.g., "#{issue_number} → $ISSUE_NUMBER"]
- [dollar sign escapes count, e.g., "Escaped 5 literal dollar signs"]

Next steps:
1. Verify the converted prompt: cat ~/.codex/prompts/$COMMAND_NAME.md
2. Test in Codex by running: /$COMMAND_NAME [arguments]
3. Check that placeholders are correctly substituted
```

## Error Handling

**Handle these errors gracefully:**

1. **Missing COMMAND_NAME parameter**:
   ```
   ❌ Error: COMMAND_NAME parameter is required

   Usage: /convert-custom-slash-for-codex COMMAND_NAME=<name>
   Example: /convert-custom-slash-for-codex COMMAND_NAME=fix-issue
   ```

2. **Source file not found**:
   ```
   ❌ Error: Command file not found: .claude/commands/$COMMAND_NAME.md

   Available commands:
   - fix-issue
   - fix-linter
   - fix-pr-review
   - create-github-issue

   Usage: /convert-custom-slash-for-codex COMMAND_NAME=<name>
   ```

3. **Cannot create output directory**:
   ```
   ❌ Error: Cannot create directory ~/.codex/prompts/
   Reason: [error message]

   Please check directory permissions and try again.
   ```

4. **Cannot write output file**:
   ```
   ❌ Error: Cannot write to ~/.codex/prompts/$COMMAND_NAME.md
   Reason: [error message]

   Please check file permissions and try again.
   ```

## Example Conversion

### Input: `.claude/commands/fix-issue.md`

```markdown
# Fix Issue #{issue_number}

You are tasked with fixing GitHub issue #{issue_number}...

Use `gh issue view {issue_number}` to fetch issue details.
```

### Output: `~/.codex/prompts/fix-issue.md`

```markdown
---
description: Fix Issue $ISSUE_NUMBER
argument-hint: ISSUE_NUMBER=#123
---
# Fix Issue $ISSUE_NUMBER

You are tasked with fixing GitHub issue $ISSUE_NUMBER...

Use `gh issue view $ISSUE_NUMBER` to fetch issue details.
```

## Testing

After conversion, manually verify:

1. **File exists**: `ls -la ~/.codex/prompts/$COMMAND_NAME.md`
2. **Content is correct**: `cat ~/.codex/prompts/$COMMAND_NAME.md`
3. **YAML frontmatter is valid**: Check description and argument-hint fields
4. **Placeholders converted**: No remaining `#{param}` or `{param}` patterns
5. **Dollar signs handled**: Literal `$` escaped to `$$`, placeholders preserved
6. **Markdown intact**: Headings, lists, code blocks properly formatted

## Limitations

**Known limitations:**
- Does not validate YAML frontmatter syntax (assumes correct format)
- Does not test the converted prompt in Codex (manual testing required)
- Parameter detection is pattern-based (may miss unusual parameter formats)
- Does not handle nested parameter references (e.g., `{{nested}}`)

## Additional Notes

- This command is read-only for source files (never modifies `.claude/commands/`)
- Overwrites existing files in `~/.codex/prompts/` without warning
- Designed for one-time conversions (not for continuous sync)
- Manual review of converted prompts is recommended before use
