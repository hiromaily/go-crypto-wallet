// Package application contains the application layer of the cryptocurrency wallet system.
//
// The application layer is organized into two main components:
//
//   - usecase/: Use case layer - Application-specific business rules and orchestration
//     Use cases coordinate domain services and infrastructure to fulfill business requirements.
//
// # Architecture
//
// The application layer follows Clean Architecture principles:
//
//   - Use cases depend on domain layer (pkg/domain/)
//   - Use cases use infrastructure through domain interfaces
//   - Use cases are organized by wallet type (watch, keygen, sign) and coin (btc, eth, xrp)
//
// # Dependency Direction
//
//	Application Layer (pkg/application/*) → Domain Layer (pkg/domain/*) ← Infrastructure Layer (pkg/infrastructure/*)
package application
