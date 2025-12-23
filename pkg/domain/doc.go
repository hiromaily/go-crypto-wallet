// Package domain contains pure business logic and domain models for the cryptocurrency wallet system.
//
// This package is the core of the application following Clean Architecture principles.
// It contains domain entities, value objects, business rules, and domain services that are
// independent of any infrastructure concerns.
//
// # Architecture Principles
//
// The domain layer:
//   - Has ZERO dependencies on infrastructure (no database, no API clients, no file I/O)
//   - Defines its own interfaces that infrastructure must implement (Dependency Inversion)
//   - Contains pure business logic that can be tested without mocks
//   - Is the most stable layer - changes here affect all layers
//
// # Dependency Direction
//
//	Application Layer (pkg/wallet/service/*) → Domain Layer (pkg/domain/*) ← Infrastructure Layer (pkg/repository/*)
//
// The domain layer sits at the center and both application and infrastructure layers depend on it,
// never the reverse.
//
// # Package Organization
//
//   - wallet/    - Wallet entities and types (watch-only, keygen, sign)
//   - account/   - Account entities and types (client, deposit, payment, stored, auth)
//   - transaction/ - Transaction entities, types, and business rules
//   - key/       - Key generation entities and value objects
//   - multisig/  - Multisig address entities and business rules
//   - coin/      - Cryptocurrency type definitions (BTC, ETH, XRP, etc.)
//
// # Domain Concepts
//
// Value Objects: Immutable objects defined by their values (AccountType, TxType, CoinTypeCode)
// Entities: Objects with unique identity and lifecycle (Transaction, Account, Key)
// Domain Services: Stateless services containing business logic that doesn't belong to a single entity
// Repository Interfaces: Defined by domain, implemented by infrastructure
//
// # Testing
//
// Domain logic should be tested without any infrastructure dependencies.
// Use pure unit tests without mocks wherever possible.
package domain
