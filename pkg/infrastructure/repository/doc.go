// Package repository provides repository implementations for data persistence.
//
// This package contains:
//   - cold/: Repository implementations for cold wallets (keygen, sign)
//   - watch/: Repository implementations for watch wallets (online)
//
// Repositories are responsible for:
//   - Implementing domain repository interfaces (when defined)
//   - CRUD operations on database tables
//   - Converting between domain entities and database models
//   - Handling database errors
//
// Repositories do NOT:
//   - Contain business logic
//   - Validate business rules (use domain validators before calling repositories)
//   - Make domain decisions
//
// Repository implementations use infrastructure/database/ for database access
// and should convert between domain entities and infrastructure database models.
package repository
