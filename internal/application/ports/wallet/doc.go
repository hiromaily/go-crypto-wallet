// Package wallet defines interfaces for wallet key generation operations.
//
// # Overview
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer (internal/infrastructure/wallet/key/).
//
// # Interfaces
//
//   - Generator: HD wallet key generation interface
//   - GeneratorFactory: Factory for creating Generator implementations
//
// # Supported Key Types
//
// The Generator interface supports multiple BIP standards:
//   - BIP44: Legacy addresses (P2PKH)
//   - BIP49: SegWit-compatible addresses (P2SH-P2WPKH)
//   - BIP84: Native SegWit addresses (P2WPKH)
//   - BIP86: Taproot addresses (P2TR)
//   - MuSig2: Multi-signature Taproot addresses
//
// # Usage
//
// Use cases use the factory to create appropriate generators:
//
//	type myUseCase struct {
//	    factory wallet.GeneratorFactory
//	}
//
//	func (u *myUseCase) GenerateKeys(keyType domainKey.KeyType) error {
//	    generator, err := u.factory.CreateGenerator(keyType, coinType, conf)
//	    if err != nil {
//	        return err
//	    }
//	    keys, err := generator.CreateKey(seed, accountType, 0, 10)
//	    // ...
//	}
//
// # Related Packages
//
//   - internal/infrastructure/wallet/key/: Generator implementations
//   - internal/domain/key/: Key types and wallet key entities
package wallet
