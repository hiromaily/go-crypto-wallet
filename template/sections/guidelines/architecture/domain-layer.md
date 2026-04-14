### Domain Layer Guidelines

The `internal/domain/` package contains pure business logic with **ZERO infrastructure dependencies**.

**Key Principles:**

- Domain layer has NO dependencies on infrastructure (no database, no API clients, no file I/O)
- Domain defines interfaces; infrastructure implements them (Dependency Inversion Principle)
- All domain logic must be testable without mocks (pure functions preferred)
- Domain is the most stable layer - changes here affect all other layers

**Domain Layer Structure:**

- **Types & Value Objects**: Immutable objects defined by values (AccountType, TxType, CoinTypeCode)
- **Entities**: Objects with unique identity and lifecycle (not yet fully implemented)
- **Validators**: Business rule validation functions
- **Domain Services**: Stateless services with business logic

**Important:**

- When adding new business logic, first consider if it belongs in the domain layer
- Use domain validators for input validation before infrastructure operations
- Business rules should be in domain, not scattered across services
