### Adding New Cryptocurrency Support

When adding support for a new cryptocurrency, follow these steps:

#### 1. Domain Layer

Add cryptocurrency-specific types and validators:

- Add coin type to `internal/domain/coin/`
- Add validators for addresses, amounts, transactions
- Define domain interfaces for the new cryptocurrency

#### 2. Infrastructure Layer

Implement API clients and repositories:

- Create API client in `internal/infrastructure/api/{coin}/`
- Implement repository interfaces in `internal/infrastructure/repository/`
- Add storage implementations if needed

#### 3. Application Layer

Create use cases for each wallet type:

- Key generation use cases in `internal/application/usecase/keygen/{coin}/`
- Signing use cases in `internal/application/usecase/sign/{coin}/`
- Watch wallet use cases in `internal/application/usecase/watch/{coin}/`

#### 4. Interface Adapters Layer

Implement CLI commands and wallet adapters:

- Add CLI commands in `internal/interface-adapters/cli/{wallet-type}/`
- Implement wallet adapter in `internal/interface-adapters/wallet/{coin}/`

#### 5. Dependency Injection

Wire up dependencies:

- Add DI container setup in `internal/di/`
- Register new cryptocurrency services

#### 6. Configuration

Add configuration support:

- Update configuration files in `config/wallet/`
- Add coin-specific settings

#### 7. Testing

Add tests for all layers:

- Domain layer: Pure unit tests
- Application layer: Use case tests with mocked infrastructure
- Infrastructure layer: Unit tests + integration tests
- Interface adapters: Command tests with mocked use cases
