# Key Generator Interface Design

This document proposes an interface-based design for key generation that allows for easy extension to support new key types (Taproot, MuSig2, etc.) as they evolve.

## Current State

### Current Implementation

```go
// internal/infrastructure/wallet/key/key.go
type Generator interface {
    CreateKey(seed []byte, actType domainAccount.AccountType, idxFrom, count uint32) ([]domainKey.WalletKey, error)
}

// internal/infrastructure/wallet/key/hd_wallet.go
type HDKey struct {
    purpose      PurposeType  // Fixed to BIP44/BIP49
    coinType     domainCoin.CoinType
    coinTypeCode domainCoin.CoinTypeCode
    conf         *chaincfg.Params
}
```

### Current Limitations

1. **Fixed Purpose Type**: `PurposeTypeBIP44` is hardcoded in `newKeyGenerator()`
2. **Single Implementation**: Only `HDKey` implements `Generator`
3. **No Key Type Selection**: Cannot choose key type (BIP44, BIP49, BIP86, etc.) dynamically
4. **Tight Coupling**: Key generation logic is tightly coupled to specific BIP standards
5. **Difficult to Extend**: Adding new key types (Taproot, MuSig2) requires modifying existing code

---

## Proposed Design

### 1. Key Type Definition

Define key types as an enum/type:

```go
// internal/domain/key/types.go
package key

// KeyType represents the type of key generation standard
type KeyType string

const (
    // KeyTypeBIP44 represents BIP44 (Legacy P2PKH)
    KeyTypeBIP44 KeyType = "bip44"
    
    // KeyTypeBIP49 represents BIP49 (P2SH-SegWit)
    KeyTypeBIP49 KeyType = "bip49"
    
    // KeyTypeBIP84 represents BIP84 (Native SegWit P2WPKH)
    KeyTypeBIP84 KeyType = "bip84"
    
    // KeyTypeBIP86 represents BIP86 (Taproot)
    KeyTypeBIP86 KeyType = "bip86"
    
    // KeyTypeMuSig2 represents MuSig2 aggregated keys
    KeyTypeMuSig2 KeyType = "musig2"
)

// String returns the string representation of the key type
func (k KeyType) String() string {
    return string(k)
}

// Purpose returns the BIP purpose number for the key type
func (k KeyType) Purpose() uint32 {
    switch k {
    case KeyTypeBIP44:
        return 44
    case KeyTypeBIP49:
        return 49
    case KeyTypeBIP84:
        return 84
    case KeyTypeBIP86:
        return 86
    default:
        return 44 // Default to BIP44
    }
}

// Validate validates the key type
func (k KeyType) Validate() error {
    switch k {
    case KeyTypeBIP44, KeyTypeBIP49, KeyTypeBIP84, KeyTypeBIP86, KeyTypeMuSig2:
        return nil
    default:
        return fmt.Errorf("invalid key type: %s", k)
    }
}
```

### 2. Enhanced Key Generator Interface

Extend the interface to support key type and additional metadata:

```go
// internal/infrastructure/wallet/key/key.go
package key

import (
    domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
    domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// Generator is the interface for key generation
type Generator interface {
    // KeyType returns the key type this generator supports
    KeyType() domainKey.KeyType
    
    // CreateKey creates keys based on the seed and account type
    CreateKey(
        seed []byte,
        accountType domainAccount.AccountType,
        idxFrom, count uint32,
    ) ([]domainKey.WalletKey, error)
    
    // SupportsAddressType checks if this generator supports the given address type
    SupportsAddressType(addrType address.AddrType) bool
    
    // GetDerivationPath returns the derivation path for the given account and index
    GetDerivationPath(accountType domainAccount.AccountType, index uint32) string
}

// GeneratorFactory creates a Generator based on key type
type GeneratorFactory interface {
    CreateGenerator(keyType domainKey.KeyType, coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) (Generator, error)
}
```

### 3. Key Type-Specific Implementations

Create separate implementations for each key type:

```go
// internal/infrastructure/wallet/key/bip44_generator.go
package key

type BIP44Generator struct {
    coinType     domainCoin.CoinType
    coinTypeCode domainCoin.CoinTypeCode
    conf         *chaincfg.Params
}

func NewBIP44Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP44Generator {
    return &BIP44Generator{
        coinType:     domainCoin.GetCoinType(coinTypeCode, conf),
        coinTypeCode: coinTypeCode,
        conf:         conf,
    }
}

func (g *BIP44Generator) KeyType() domainKey.KeyType {
    return domainKey.KeyTypeBIP44
}

func (g *BIP44Generator) CreateKey(
    seed []byte,
    accountType domainAccount.AccountType,
    idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
    // BIP44 implementation
    // Derivation path: m/44'/coin_type'/account'/0/index
    // ...
}

func (g *BIP44Generator) SupportsAddressType(addrType address.AddrType) bool {
    return addrType == address.AddrTypeLegacy
}

func (g *BIP44Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
    return fmt.Sprintf("m/44'/%d'/%d'/0/%d", g.coinType.Uint32(), accountType.Uint32(), index)
}
```

```go
// internal/infrastructure/wallet/key/bip86_generator.go
package key

type BIP86Generator struct {
    coinType     domainCoin.CoinType
    coinTypeCode domainCoin.CoinTypeCode
    conf         *chaincfg.Params
}

func NewBIP86Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP86Generator {
    return &BIP86Generator{
        coinType:     domainCoin.GetCoinType(coinTypeCode, conf),
        coinTypeCode: coinTypeCode,
        conf:         conf,
    }
}

func (g *BIP86Generator) KeyType() domainKey.KeyType {
    return domainKey.KeyTypeBIP86
}

func (g *BIP86Generator) CreateKey(
    seed []byte,
    accountType domainAccount.AccountType,
    idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
    // BIP86 implementation
    // Derivation path: m/86'/coin_type'/account'/0/index
    // Generate Taproot addresses (bc1p...)
    // ...
}

func (g *BIP86Generator) SupportsAddressType(addrType address.AddrType) bool {
    // Taproot addresses
    return addrType == address.AddrTypeTaproot
}

func (g *BIP86Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
    return fmt.Sprintf("m/86'/%d'/%d'/0/%d", g.coinType.Uint32(), accountType.Uint32(), index)
}
```

### 4. Generator Factory

Create a factory to instantiate the appropriate generator:

```go
// internal/infrastructure/wallet/key/factory.go
package key

import (
    "fmt"
    
    domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
    domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
    "github.com/btcsuite/btcd/chaincfg"
)

// Factory creates key generators based on key type
type Factory struct{}

func NewFactory() *Factory {
    return &Factory{}
}

// CreateGenerator creates a generator for the specified key type
func (f *Factory) CreateGenerator(
    keyType domainKey.KeyType,
    coinTypeCode domainCoin.CoinTypeCode,
    conf *chaincfg.Params,
) (Generator, error) {
    if err := keyType.Validate(); err != nil {
        return nil, fmt.Errorf("invalid key type: %w", err)
    }
    
    switch keyType {
    case domainKey.KeyTypeBIP44:
        return NewBIP44Generator(coinTypeCode, conf), nil
    case domainKey.KeyTypeBIP49:
        return NewBIP49Generator(coinTypeCode, conf), nil
    case domainKey.KeyTypeBIP84:
        return NewBIP84Generator(coinTypeCode, conf), nil
    case domainKey.KeyTypeBIP86:
        return NewBIP86Generator(coinTypeCode, conf), nil
    case domainKey.KeyTypeMuSig2:
        return NewMuSig2Generator(coinTypeCode, conf), nil
    default:
        return nil, fmt.Errorf("unsupported key type: %s", keyType)
    }
}

// CreateGeneratorFromConfig creates a generator based on configuration
func (f *Factory) CreateGeneratorFromConfig(
    keyType domainKey.KeyType,
    coinTypeCode domainCoin.CoinTypeCode,
    networkType string,
) (Generator, error) {
    conf := getChainConfig(coinTypeCode, networkType)
    return f.CreateGenerator(keyType, coinTypeCode, conf)
}
```

### 5. Configuration Support

Add key type to configuration:

```go
// pkg/config/wallet_struct.go
type WalletRoot struct {
    // ... existing fields
    KeyType      domainKey.KeyType `toml:"key_type" mapstructure:"key_type" validate:"oneof=bip44 bip49 bip84 bip86 musig2"`
    AddressType  address.AddrType   `toml:"address_type" mapstructure:"address_type" validate:"oneof=p2sh-segwit bech32 bch-cashaddr taproot"`
    // ...
}
```

### 6. Updated DI Container

Update the dependency injection container to use the factory:

```go
// internal/di/container.go
func (c *container) newKeyGenerator() key.Generator {
    var chainConf *chaincfg.Params
    switch {
    case domainCoin.IsBTCGroup(c.conf.CoinTypeCode):
        chainConf = c.newBTC().GetChainConf()
    case domainCoin.IsETHGroup(c.conf.CoinTypeCode):
        chainConf = c.newETH().GetChainConf()
    case c.conf.CoinTypeCode == domainCoin.XRP:
        chainConf = c.newXRP().GetChainConf()
    default:
        panic(fmt.Sprintf("coinType[%s] is not implemented yet.", c.conf.CoinTypeCode))
    }
    
    // Use factory to create generator based on key type
    factory := key.NewFactory()
    keyType := c.getKeyType() // Get from config or default to BIP44
    generator, err := factory.CreateGenerator(keyType, c.conf.CoinTypeCode, chainConf)
    if err != nil {
        panic(fmt.Sprintf("failed to create key generator: %v", err))
    }
    
    return generator
}

func (c *container) getKeyType() domainKey.KeyType {
    // Get from config if available, otherwise default to BIP44
    if c.conf.KeyType != "" {
        return c.conf.KeyType
    }
    return domainKey.KeyTypeBIP44
}
```

### 7. Backward Compatibility

Maintain backward compatibility by making `HDKey` implement the new interface:

```go
// internal/infrastructure/wallet/key/hd_wallet.go
// HDKey implements Generator interface (backward compatibility)
func (k *HDKey) KeyType() domainKey.KeyType {
    switch k.purpose {
    case PurposeTypeBIP44:
        return domainKey.KeyTypeBIP44
    case PurposeTypeBIP49:
        return domainKey.KeyTypeBIP49
    default:
        return domainKey.KeyTypeBIP44
    }
}

func (k *HDKey) SupportsAddressType(addrType address.AddrType) bool {
    switch k.purpose {
    case PurposeTypeBIP44:
        return addrType == address.AddrTypeLegacy
    case PurposeTypeBIP49:
        return addrType == address.AddrTypeP2shSegwit
    default:
        return false
    }
}

func (k *HDKey) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
    return fmt.Sprintf("m/%d'/%d'/%d'/0/%d", 
        k.purpose.Uint32(), 
        k.coinType.Uint32(), 
        accountType.Uint32(), 
        index)
}
```

---

## Benefits

### 1. **Extensibility**
- Easy to add new key types (BIP86, MuSig2, etc.)
- No need to modify existing code
- Clear separation of concerns

### 2. **Testability**
- Each generator can be tested independently
- Mock generators for testing
- Clear interface contracts

### 3. **Flexibility**
- Choose key type at runtime based on configuration
- Support multiple key types simultaneously
- Easy to switch between key types

### 4. **Maintainability**
- Clear code organization
- Each key type has its own implementation
- Easier to understand and modify

### 5. **Future-Proof**
- Ready for new Bitcoin improvements
- Easy to add quantum-resistant key types
- Supports evolving standards

---

## Migration Strategy

### Phase 1: Interface Definition (1-2 weeks)

1. Define `KeyType` in domain layer
2. Extend `Generator` interface
3. Create factory interface

### Phase 2: Implementation (2-4 weeks)

1. Create `BIP44Generator` (refactor from `HDKey`)
2. Create `BIP49Generator`
3. Create `BIP86Generator` (new)
4. Update factory implementation

### Phase 3: Integration (1-2 weeks)

1. Update DI container
2. Add configuration support
3. Update use cases to use new interface

### Phase 4: Testing & Documentation (1-2 weeks)

1. Unit tests for each generator
2. Integration tests
3. Update documentation

### Phase 5: Deprecation (Optional)

1. Mark old `HDKey` as deprecated
2. Migrate existing code to use new generators
3. Remove old implementation after migration

---

## Example Usage

### Configuration

```toml
# data/config/btc_watch.toml
[wallet]
coin_type = "btc"
key_type = "bip86"  # Use Taproot
address_type = "taproot"
```

### Code Usage

```go
// Create generator factory
factory := key.NewFactory()

// Create BIP86 generator for Taproot
generator, err := factory.CreateGenerator(
    domainKey.KeyTypeBIP86,
    domainCoin.BTC,
    &chaincfg.MainNetParams,
)
if err != nil {
    return err
}

// Generate keys
walletKeys, err := generator.CreateKey(seed, accountType, 0, 10)
if err != nil {
    return err
}

// Get derivation path
path := generator.GetDerivationPath(accountType, 0)
// Output: m/86'/0'/0'/0/0
```

---

## Implementation Checklist

- [ ] Define `KeyType` in domain layer
- [ ] Extend `Generator` interface
- [ ] Create `BIP44Generator`
- [ ] Create `BIP49Generator`
- [ ] Create `BIP84Generator`
- [ ] Create `BIP86Generator` (Taproot)
- [ ] Create `MuSig2Generator` (future)
- [ ] Create `Factory`
- [ ] Update DI container
- [ ] Add configuration support
- [ ] Update use cases
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Update documentation
- [ ] Migration guide

---

## Summary

This design provides:

1. **Interface-based architecture** for key generation
2. **Type-specific implementations** for each key standard
3. **Factory pattern** for creating generators
4. **Configuration-driven** key type selection
5. **Easy extensibility** for future key types
6. **Backward compatibility** with existing code

This approach allows the system to evolve with Bitcoin standards while maintaining clean, testable, and maintainable code.

