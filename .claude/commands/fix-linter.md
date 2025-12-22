# Fix Linter Command

You are tasked with fixing lint errors that occurred when running the `make lint-fix` command.

## Guidelines

1. **Analyze the errors**: First, understand the types and severity of lint errors present
2. **Prioritize fixes**: Focus on the most critical errors first:
   - Syntax errors and breaking issues
   - Security vulnerabilities
   - Type errors
   - Style and formatting issues (lowest priority)
3. **Step-by-step approach**: Fix errors incrementally and verify each change
4. **Batch similar errors**: Group and fix similar error types together for efficiency
5. **Preserve functionality**: Ensure all fixes maintain the original code behavior
6. **Use auto-fix when possible**: Leverage linter's automatic fixing capabilities where appropriate

## Process

- Start by showing a summary of the error types and their counts
- Propose a fix order based on severity and impact
- Fix the most critical errors first
- After each significant fix, explain what was changed and why
- If there are many errors (>20), break the work into logical chunks
- Verify the fixes don't introduce new issues

## Output Format

For each fix iteration:

1. Describe the error category being addressed
2. Show the specific changes made
3. Explain the reasoning behind the fix
4. Note any remaining errors to address in subsequent steps

Remember: Quality over speed. It's better to fix errors correctly in stages than to rush through all of them at once.
