### Import Organization

Organize imports in the following order:

1. Standard library packages
2. Third-party packages
3. Local packages (this project)

Example:

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third-party packages
    "github.com/btcsuite/btcd/btcutil"
    "github.com/pkg/errors"

    // Local packages
    "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
    "github.com/hiromaily/go-crypto-wallet/pkg/logger"
)
```

The `goimports` tool (via `make go-fmt`) will automatically organize imports in this order.
