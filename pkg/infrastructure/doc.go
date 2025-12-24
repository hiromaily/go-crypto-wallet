// Package infrastructure provides implementations for external system integrations
// following Clean Architecture principles.
//
// The infrastructure layer is organized by technical concern and implements interfaces
// defined in the domain layer (Dependency Inversion Principle).
//
// # Package Organization
//
// The infrastructure layer is organized into the following subdirectories:
//
//   - database/: Database connections and query execution (MySQL, sqlc-generated code)
//   - repository/: Repository implementations for data persistence (cold wallet, watch wallet)
//   - api/: External blockchain API clients (Bitcoin, Ethereum, Ripple)
//   - storage/: File-based storage implementations (transaction files, address files)
//   - network/: Network connection management (WebSocket, gRPC)
//
// # Dependency Direction
//
// Infrastructure implements domain interfaces:
//
//	Application Layer (wallet/service) → Domain Layer (domain/*) ← Infrastructure Layer (infrastructure/*)
//
// The domain layer defines interfaces, and the infrastructure layer provides concrete
// implementations. This allows the application layer to depend on stable domain abstractions
// rather than volatile infrastructure details.
//
// # Infrastructure vs. Shared Utilities
//
// Infrastructure components represent external systems and technical implementations:
//   - Database connections and repositories
//   - External API clients (blockchain nodes)
//   - File I/O for transaction/address storage
//   - Network connections (WebSocket, gRPC)
//
// Shared utilities remain in pkg/ and are used across all layers:
//   - uuid/: UUID generation
//   - logger/: Logging utilities
//   - config/: Configuration management
//   - converter/: Type conversion utilities
//   - serial/: Serialization utilities
//   - debug/: Debug utilities
//   - di/: Dependency injection container
//   - testutil/: Test utilities
//   - wallet/key/: Cryptographic key generation (HD wallet, seeds)
//   - contract/: Smart contract interaction utilities
//
// # Infrastructure Component Guidelines
//
// Infrastructure components should:
//   - Implement domain interfaces (when defined)
//   - Contain NO business logic (only technical implementation)
//   - Be easily replaceable and mockable for testing
//   - Handle external system communication and error handling
//   - Convert between domain entities and external system formats
//
// # Database Infrastructure (database/)
//
// Responsibilities:
//   - Database connection management
//   - Query execution via sqlc-generated code
//   - Transaction management
//
// Does NOT:
//   - Contain business logic
//   - Validate business rules
//   - Make domain decisions
//
// # Repository Implementations (repository/)
//
// Responsibilities:
//   - Implement domain repository interfaces (if defined)
//   - CRUD operations on database tables
//   - Convert between domain entities and database models
//
// Does NOT:
//   - Contain business logic
//   - Validate business rules (domain validators should be used)
//   - Make domain decisions
//
// # API Clients (api/)
//
// Responsibilities:
//   - Communicate with external blockchain APIs (Bitcoin Core, Ethereum, Ripple)
//   - Handle network errors and retries
//   - Serialize/deserialize API requests and responses
//
// Does NOT:
//   - Contain business logic
//   - Validate transactions (domain validators should be used)
//   - Make domain decisions
//
// # File Storage (storage/)
//
// Responsibilities:
//   - Read/write files for transaction and address data
//   - Manage file paths and formats
//   - Handle file I/O errors
//
// Does NOT:
//   - Contain business logic
//   - Validate file content (beyond format validation)
//   - Make domain decisions
//
// # Network Infrastructure (network/)
//
// Responsibilities:
//   - Establish and manage network connections (WebSocket, gRPC)
//   - Handle connection lifecycle and errors
//   - Provide connection objects to API clients
//
// Does NOT:
//   - Contain business logic
//   - Make routing decisions
//   - Handle application-level protocols
//
// # Testing Strategy
//
// Infrastructure components should be tested in two ways:
//
//  1. Unit Tests: Test infrastructure in isolation with mocked external dependencies
//  2. Integration Tests: Test infrastructure with real external dependencies
//
// Infrastructure components should be easily mockable to allow application layer
// testing without real external dependencies.
//
// # Adding New Infrastructure Components
//
// When adding new infrastructure components:
//
//  1. Determine the appropriate subdirectory (database, repository, api, storage, network)
//  2. Implement domain interfaces if they exist
//  3. Keep business logic in domain layer, technical implementation in infrastructure
//  4. Add package documentation explaining the component's purpose
//  5. Create mock implementations for testing
//  6. Update dependency injection container to wire the component
//
// # References
//
//   - Clean Architecture: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
//   - Dependency Inversion Principle: https://en.wikipedia.org/wiki/Dependency_inversion_principle
//   - Project's AGENTS.md for coding standards
//   - Project's ISSUE_SEPARATE_INFRASTRUCTURE_LAYER.md for detailed refactoring plan
package infrastructure
