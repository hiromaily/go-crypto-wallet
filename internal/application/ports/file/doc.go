// Package file defines interfaces for file system I/O operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/storage/file/).
//
// # Interfaces
//
//   - TransactionFileRepositorier: Transaction file read/write operations
//   - AddressFileRepositorier: Address CSV file operations
//   - DescriptorFileWriter: Descriptor file write operations
//
// # Types
//
//   - FileName: Parsed transaction file name components
//
// # Capabilities
//
// TransactionFileRepositorier provides:
//   - File path creation and validation
//   - Hex transaction read/write (WriteFile, ReadFile)
//   - Multi-line data read/write (WriteFileSlice, ReadFileSlice)
//   - PSBT file support (ReadPSBTFile, WritePSBTFile) - BIP174
//
// AddressFileRepositorier provides:
//   - Address CSV file import
//   - File path validation
//
// # Usage
//
// Use cases depend on these interfaces, not concrete implementations:
//
//	type myUseCase struct {
//	    txFile   file.TransactionFileRepositorier
//	    addrFile file.AddressFileRepositorier
//	}
//
// # Related Packages
//
//   - internal/infrastructure/storage/file/transaction/: Transaction file implementation
//   - internal/infrastructure/storage/file/address/: Address file implementation
package file
