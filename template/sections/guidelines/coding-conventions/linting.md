### Linting and Formatting

- Follow `golangci-lint` configuration (`.golangci.yml`)
- Format code with `make go-fmt` (uses `gofumpt` and `goimports` via golangci-lint)
  - Import order: standard → third-party → local
- Use `make go-lint` to run linting and formatting together (executes lint checks and format fixes)
- Maintain consistent naming conventions (lowercase package names, exported functions start with uppercase)
