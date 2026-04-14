### Linter-Specific Guidelines

**unused-receiver:**

- For `unused-receiver` lint errors: **Remove the receiver entirely** instead of renaming it to `_`
- Renaming to `_` will only cause the same error to appear for other receivers
- Convert the method to a function from the start

**errcheck:**

- Never ignore errors
- If an error must be ignored, add a comment explaining why
- Use `_ = err` with a comment for intentionally ignored errors

**gofmt / gofumpt:**

- Always run `make go-fmt` before committing
- Consistent formatting helps with code reviews

**goimports:**

- Import organization is handled by `goimports`
- Run `make go-fmt` to organize imports automatically
