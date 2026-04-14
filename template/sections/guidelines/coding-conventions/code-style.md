### Code Style

**Function Length:**

- Keep functions short and focused
- Prefer small, composable functions over large, complex ones
- If a function exceeds 50 lines, consider refactoring

**Error Handling:**

- Always check errors
- Wrap errors with context using `fmt.Errorf` with `%w`
- Return errors early (early return pattern)

**Comments:**

- Add godoc comments to all exported functions, methods, types, and constants
- Keep comments up-to-date with code changes
- Use complete sentences in comments
- Explain "why" rather than "what" in implementation comments

**Variable Naming:**

- Use short names for short-lived variables (e.g., `i` for loop index, `err` for errors)
- Use descriptive names for long-lived variables
- Avoid single-letter names except for:
  - Loop indices (`i`, `j`, `k`)
  - Short-lived variables in small scopes
  - Receivers (use consistent short names like `r` for `Repository`)
