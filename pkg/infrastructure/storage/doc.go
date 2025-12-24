// Package storage provides file-based storage implementations.
//
// This package contains:
//   - file/: File I/O for transaction and address data exchange
//
// File storage is used to transfer data between wallet types:
//   - Transaction hex data (unsigned/signed) between keygen, sign, and watch wallets
//   - Address data exported from keygen wallet and imported to watch wallet
//
// File storage is responsible for:
//   - Reading and writing transaction files
//   - Reading and writing address CSV files
//   - Creating structured filenames with metadata
//   - Handling file I/O errors
//
// File storage does NOT:
//   - Contain business logic
//   - Validate transaction data (use domain validators)
//   - Make decisions about transaction flow
//
// File storage implements domain storage interfaces (if defined) to allow
// the application layer to work with storage abstractions.
package storage
