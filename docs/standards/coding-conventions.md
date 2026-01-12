# Coding Conventions

Project coding standards for go-crypto-wallet.

## Go Code

### Verification Commands

```bash
make go-lint      # Lint and format
make check-build  # Verify build
make gotest       # Run tests
make tidy         # Clean dependencies
```

### Import Order

1. Standard library
2. Third-party packages
3. Local packages

```go
import (
    "context"
    "fmt"

    "github.com/btcsuite/btcd/btcutil"

    "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Package | lowercase, no underscores | `account`, `transaction` |
| Exported | UpperCamelCase | `GetAccountKey` |
| Unexported | lowerCamelCase | `calculateFee` |
| Interface | Behavior + "er" | `Validator`, `Reader` |

### Error Handling

```go
result, err := service.DoSomething()
if err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}
```

## TypeScript/JavaScript

### Verification Commands

```bash
cd apps/{app-name}
npm run lint
npm run format
npm run build
npm test
```

## Shell Scripts

```bash
make shfmt  # Format shell scripts
```

## Detailed Guidelines

See [docs/ai-agents/guidelines/coding-standards.md](../ai-agents/guidelines/coding-standards.md) for full details.
