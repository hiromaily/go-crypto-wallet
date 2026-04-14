### Application Ports Layer (Interface Definitions)

The `internal/application/ports/` package contains **interface definitions** (abstractions) for infrastructure dependencies, following the Dependency Inversion Principle.

**Key Principles:**

- **Interfaces are defined in the layer that uses them** (application layer), NOT in the infrastructure layer
- **Infrastructure layer contains only implementations**, never interface definitions
- This inverts the dependency direction: Infrastructure depends on Application, not vice versa
- Ports define contracts that infrastructure implementations must fulfill
- **Interfaces MUST use application-layer DTOs**, NOT infrastructure types (Clean Architecture Dependency Rule)

**Current Port Packages:**

```text
internal/application/ports/
├── bitcoin/
│   └── interface.go              # Bitcoiner interface (Bitcoin/BCH API abstraction)
├── persistence/
│   └── repository.go             # Repository interfaces (database abstractions)
└── storage/
    └── interface.go              # TransactionFileRepositorier interface (file storage abstraction)
```

**DTOs for Port Interfaces:**

Port interfaces must use **application-layer DTOs** (defined in `internal/application/dto/`) instead of infrastructure types. This ensures the application layer has zero dependencies on infrastructure.

**DTO Location:**

- DTOs are defined in `internal/application/dto/{coin}/` (e.g., `internal/application/dto/btc/dto.go`)
- DTOs contain only domain types, standard library types, or external library types (e.g., `btcutil.Amount`)
- DTOs MUST NOT import infrastructure packages

**Example: Bitcoin API Abstraction with DTOs**

```go
// DTOs in internal/application/dto/btc/dto.go
package btc

import "github.com/btcsuite/btcd/btcutil"

type AddressInfo struct {
    Address      string
    ScriptPubKey string
    IsWitness    bool
    Labels       []string
    // ... other fields
}

type UnspentOutput struct {
    TxID          string
    Vout          uint32
    Address       string
    Amount        btcutil.Amount
    Confirmations int64
    // ... other fields
}

// Interface definition in application/ports/btc/interface.go
package btc

import btcdto "internal/application/dto/btc"

type Bitcoiner interface {
    GetBalance() (btcutil.Amount, error)
    GetAddressInfo(addr string) (*btcdto.AddressInfo, error)
    ListUnspent(confirmationNum uint64) ([]btcdto.UnspentOutput, error)
    // ... other methods using DTOs
}

// Implementation in infrastructure/api/btc/btc/bitcoin.go
package btc

import (
    portsBtc "internal/application/ports/btc"
    btcdto "internal/application/dto/btc"
)

type Bitcoin struct {
    client *rpcclient.Client
}

// Bitcoin implements portsBtc.Bitcoiner interface
// Maps infrastructure types to application DTOs
func (b *Bitcoin) GetAddressInfo(addr string) (*btcdto.AddressInfo, error) {
    // Call Bitcoin Core RPC (returns infrastructure type)
    result, err := b.client.GetAddressInfo(addr)
    if err != nil {
        return nil, err
    }
    
    // Map infrastructure type to application DTO
    return &btcdto.AddressInfo{
        Address:      result.Address,
        ScriptPubKey: result.ScriptPubKey,
        IsWitness:    result.IsWitness,
        Labels:       result.Labels,
        // ... map all fields
    }, nil
}
```

**Why Interfaces Belong in Application Layer:**

1. **Dependency Inversion Principle**: High-level modules (application) should not depend on low-level modules (infrastructure). Both should depend on abstractions.
2. **Stable Abstractions**: The application defines what it needs; infrastructure provides implementations.
3. **Testability**: Application layer can be tested with mock implementations without infrastructure dependencies.
4. **Clean Architecture**: Core business logic (application) is independent of external frameworks and libraries (infrastructure).

**Important:**

- **NEVER** define interfaces in the infrastructure layer
- **ALWAYS** define interfaces in `application/ports/` when infrastructure needs abstraction
- **NEVER** use infrastructure types in port interface method signatures (use application DTOs instead)
- **ALWAYS** create DTOs in `internal/application/dto/` for data structures returned by infrastructure
- Infrastructure packages import and implement these port interfaces
- Infrastructure implementations map between infrastructure types and application DTOs
- Use cases depend on port interfaces, not concrete implementations

**DTO Creation Guidelines:**

When creating port interfaces that return data from infrastructure:

1. **Identify infrastructure types** used in interface methods
2. **Create corresponding DTOs** in `internal/application/dto/{coin}/dto.go` (e.g., `internal/application/dto/btc/dto.go`)
3. **Update interface methods** to use DTOs instead of infrastructure types
4. **Update infrastructure implementations** to map infrastructure types → DTOs
5. **Update use cases** to work with new DTOs

**Example Pattern:**

```go
// ❌ WRONG: Interface depends on infrastructure type
type Bitcoiner interface {
    GetAddressInfo(addr string) (*btc.GetAddressInfoResult, error)  // btc is infrastructure
}

// ✅ CORRECT: Interface uses application DTO
type Bitcoiner interface {
    GetAddressInfo(addr string) (*btcdto.AddressInfo, error)  // btcdto is application DTO
}
```
