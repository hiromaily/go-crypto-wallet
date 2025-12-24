# Create GitHub Issue for AI Implementation

This command helps create a well-structured GitHub issue that is suitable for AI agents (like Cursor) to automatically implement features or fix bugs.

## How to Use

### In Cursor Chat

1. Open Cursor's chat interface (Cmd+L or Ctrl+L)
2. Type `@create-github-issue` or reference this command file
3. Describe what you want to implement or fix
4. The AI will generate a well-structured GitHub issue following the template below

### Example Usage

```
@create-github-issue I need to add a new validator function for account types in the domain layer. 
It should validate that account types follow Clean Architecture principles and be placed in 
pkg/domain/account/validator.go. Reference the existing validators in pkg/domain/transaction/ 
for the pattern.
```

The AI will then generate a complete GitHub issue with all required sections formatted as Markdown.

### Manual Usage

You can also read this file directly and use it as a guide when creating issues manually on GitHub.

## Purpose

Generate a GitHub issue in Markdown format that contains all necessary information for AI agents to understand the requirements and implement the solution automatically.

## Instructions

When creating a GitHub issue, include the following sections to ensure AI agents can implement it effectively:

### Required Sections

1. **Title**: Clear, concise description of the feature or bug
2. **Description**: Detailed explanation of what needs to be implemented
3. **Acceptance Criteria**: Specific, testable conditions that define completion
4. **Technical Requirements**: Architecture, patterns, and constraints to follow
5. **Implementation Details**: Code examples, file locations, and related code references
6. **Testing Requirements**: How to verify the implementation
7. **Related Context**: Links to related issues, PRs, or documentation

### Template Structure

Use the following template structure:

```markdown
## Description
[Clear description of what needs to be implemented]

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Technical Requirements
- Architecture: [Clean Architecture principles, layer separation, etc.]
- Patterns: [Dependency injection, interfaces, etc.]
- Constraints: [Security requirements, performance, etc.]

## Implementation Details
- Files to modify/create: [List of files]
- Related code: [Code references with line numbers]
- Dependencies: [Any new dependencies needed]

## Testing Requirements
- Unit tests: [What to test]
- Integration tests: [If applicable]
- Manual testing: [Steps to verify]

## Related Context
- Related issues: [Issue numbers]
- Documentation: [Links to relevant docs]
- Architecture guidelines: [References to AGENTS.md, etc.]
```

## Guidelines for AI-Friendly Issues

1. **Be Specific**: Include exact file paths, function names, and code references
2. **Provide Context**: Reference existing code patterns and architecture
3. **Define Boundaries**: Clearly state what should and shouldn't be changed
4. **Include Examples**: Show expected code structure or behavior
5. **Specify Constraints**: Mention security, performance, or architectural constraints
6. **Reference Documentation**: Link to relevant project documentation (AGENTS.md, etc.)

## Output Format

The issue should be formatted as valid Markdown that can be directly pasted into GitHub's issue editor.

**Note**: The generated markdown file should be saved to the `docs/issues/` directory.
