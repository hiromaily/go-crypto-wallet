# Infrastructure Layer Repository Audit Report

**Date:** 2026-01-02
**Related Issue:** #224 - Refactor infrastructure layer to align with domain-aware I/O design principles
**Sub-Issue:** #225 - Audit and Document Current State

## Executive Summary

This audit documents all repository implementations in `internal/infrastructure/repository/` that return `sqlcgen` types, interface definitions, and existing domain entities. This analysis supports the refactoring effort to align with Clean Architecture principles where repositories should return domain entities instead of database-specific types.

**Key Findings:**
- 19 repository files returning `sqlcgen` types
- 14 distinct `sqlcgen` types requiring domain entity conversion
- Repository interfaces already migrated to `application/ports/persistence/`
- 1 exemplary implementation already following the pattern (`NonceRepositorySqlc`)

---

## 1. Repository Files Returning `sqlcgen` Types

### 1.1 Watch Wallet Repository Files

#### Address Repository
**File:** `internal/infrastructure/repository/watch/address_sqlc.go`
**Struct:** `AddressRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetAll(accountType domainAccount.AccountType) ([]*sqlcgen.Address, error)`
- `GetOneUnAllocated(accountType domainAccount.AccountType) (*sqlcgen.Address, error)`
- `InsertBulk(ctx context.Context, items []*sqlcgen.Address) error`

**sqlcgen types used:** `sqlcgen.Address`

**sqlcgen.Address structure:**
```go
- ID              int64
- Coin            AddressCoin
- Account         AddressAccount
- WalletAddress   string
- IsAllocated     bool
- UpdatedAt       sql.NullTime
```

---

#### Transaction Repository
**File:** `internal/infrastructure/repository/watch/tx_sqlc.go`
**Struct:** `TxRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(id int64) (*sqlcgen.Tx, error)`
- `Update(txItem *sqlcgen.Tx) (int64, error)`

**sqlcgen types used:** `sqlcgen.Tx`

**Special features:**
- Implements `WithTx(tx *sql.Tx) portsPersistence.TxRepositorier` for transaction support

**sqlcgen.Tx structure:**
```go
- ID          int64
- Coin        TxCoin
- Action      TxAction
- UpdatedAt   sql.NullTime
```

---

#### BTC Transaction Repository
**File:** `internal/infrastructure/repository/watch/btc_tx_sqlc.go`
**Struct:** `BTCTxRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(id int64) (*sqlcgen.BtcTx, error)`
- `InsertUnsignedTx(actionType domainTx.ActionType, txItem *sqlcgen.BtcTx) (int64, error)`
- `Update(txItem *sqlcgen.BtcTx) (int64, error)`

**sqlcgen types used:** `sqlcgen.BtcTx`

**sqlcgen.BtcTx structure:**
```go
- ID                  int64
- Coin                BtcTxCoin
- Action              BtcTxAction
- UnsignedHexTx       string
- SignedHexTx         string
- SentHashTx          string
- TotalInputAmount    string
- TotalOutputAmount   string
- Fee                 string
- CurrentTxType       int8
- UnsignedUpdatedAt   sql.NullTime
- SentUpdatedAt       sql.NullTime
```

---

#### BTC Transaction Input Repository
**File:** `internal/infrastructure/repository/watch/btc_tx_input_sqlc.go`
**Struct:** `TxInputRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(id int64) (*sqlcgen.BtcTxInput, error)`
- `GetAllByTxID(id int64) ([]*sqlcgen.BtcTxInput, error)`
- `Insert(txItem *sqlcgen.BtcTxInput) error`
- `InsertBulk(txItems []*sqlcgen.BtcTxInput) error`

**sqlcgen types used:** `sqlcgen.BtcTxInput`

**sqlcgen.BtcTxInput structure:**
```go
- ID                  int64
- TxID                int64
- InputTxid           string
- InputVout           uint32
- InputAddress        string
- InputAccount        string
- InputAmount         string
- InputConfirmations  uint64
- UpdatedAt           sql.NullTime
```

---

#### BTC Transaction Output Repository
**File:** `internal/infrastructure/repository/watch/btc_tx_output_sqlc.go`
**Struct:** `TxOutputRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(id int64) (*sqlcgen.BtcTxOutput, error)`
- `GetAllByTxID(id int64) ([]*sqlcgen.BtcTxOutput, error)`
- `Insert(txItem *sqlcgen.BtcTxOutput) error`
- `InsertBulk(txItems []*sqlcgen.BtcTxOutput) error`

**sqlcgen types used:** `sqlcgen.BtcTxOutput`

**sqlcgen.BtcTxOutput structure:**
```go
- ID              int64
- TxID            int64
- OutputAddress   string
- OutputAccount   string
- OutputAmount    string
- IsChange        bool
- UpdatedAt       sql.NullTime
```

---

#### Payment Request Repository
**File:** `internal/infrastructure/repository/watch/payment_request_sqlc.go`
**Struct:** `PaymentRequestRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetAll() ([]*sqlcgen.PaymentRequest, error)`
- `GetAllByPaymentID(paymentID int64) ([]*sqlcgen.PaymentRequest, error)`
- `InsertBulk(items []*sqlcgen.PaymentRequest) error`

**sqlcgen types used:** `sqlcgen.PaymentRequest`

**Special features:**
- Implements `WithTx(tx *sql.Tx) portsPersistence.PaymentRequestRepositorier`

**sqlcgen.PaymentRequest structure:**
```go
- ID              int64
- Coin            PaymentRequestCoin
- PaymentID       sql.NullInt64
- SenderAddress   string
- SenderAccount   string
- ReceiverAddress string
- Amount          string
- IsDone          bool
- UpdatedAt       sql.NullTime
```

---

#### ETH Detail Transaction Repository
**File:** `internal/infrastructure/repository/watch/eth_detail_tx_sqlc.go`
**Struct:** `ETHDetailTXInputRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(id int64) (*sqlcgen.EthDetailTx, error)`
- `GetAllByTxID(id int64) ([]*sqlcgen.EthDetailTx, error)`
- `Insert(txItem *sqlcgen.EthDetailTx) error`
- `InsertBulk(txItems []*sqlcgen.EthDetailTx) error`

**sqlcgen types used:** `sqlcgen.EthDetailTx`

**Special features:**
- Implements `WithTx(tx *sql.Tx) portsPersistence.ETHDetailTXRepositorier`

**sqlcgen.EthDetailTx structure:**
```go
- ID                  int64
- TxID                int64
- Uuid                string
- CurrentTxType       int8
- SenderAccount       string
- SenderAddress       string
- ReceiverAccount     string
- ReceiverAddress     string
- Amount              uint64
- Fee                 uint64
- GasLimit            uint64
- Nonce               uint64
- UnsignedHexTx       string
- SignedHexTx         string
- SentHashTx          string
- UnsignedUpdatedAt   sql.NullTime
- SentUpdatedAt       sql.NullTime
```

---

#### XRP Detail Transaction Repository
**File:** `internal/infrastructure/repository/watch/xrp_detail_tx_sqlc.go`
**Struct:** `XRPDetailTxInputRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(id int64) (*sqlcgen.XrpDetailTx, error)`
- `GetAllByTxID(id int64) ([]*sqlcgen.XrpDetailTx, error)`
- `Insert(txItem *sqlcgen.XrpDetailTx) error`
- `InsertBulk(txItems []*sqlcgen.XrpDetailTx) error`

**sqlcgen types used:** `sqlcgen.XrpDetailTx`

**Special features:**
- Implements `WithTx(tx *sql.Tx) portsPersistence.XRPDetailTXRepositorier`

**sqlcgen.XrpDetailTx structure:**
```go
- ID                     int64
- TxID                   int64
- Uuid                   string
- CurrentTxType          int8
- SenderAccount          string
- SenderAddress          string
- ReceiverAccount        string
- ReceiverAddress        string
- Amount                 string
- XrpTxType              string
- Fee                    uint64
- Flags                  uint32
- LastLedgerSequence     uint32
- Sequence               uint32
- SigningPubkey          string
- TxnSignature           string
- Hash                   string
- EarliestLedgerVersion  uint64
- SignedTxID             string
- TxBlob                 string
- SentUpdatedAt          sql.NullTime
```

---

### 1.2 Cold Wallet Repository Files

#### Seed Repository
**File:** `internal/infrastructure/repository/cold/seed_sqlc.go`
**Struct:** `SeedRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(ctx context.Context) (*sqlcgen.Seed, error)`

**sqlcgen types used:** `sqlcgen.Seed`

**sqlcgen.Seed structure:**
```go
- ID          int8
- Coin        SeedCoin
- Seed        string
- UpdatedAt   sql.NullTime
```

**Security Note:** Seed contains sensitive cryptographic material

---

#### BTC Account Key Repository
**File:** `internal/infrastructure/repository/cold/account_key_sqlc.go`
**Struct:** `BTCAccountKeyRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOneMaxID(accountType domainAccount.AccountType) (*sqlcgen.BtcAccountKey, error)`
- `GetAllAddrStatus(accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus) ([]*sqlcgen.BtcAccountKey, error)`
- `GetAllMultiAddr(accountType domainAccount.AccountType, addrs []string) ([]*sqlcgen.BtcAccountKey, error)`
- `InsertBulk(items []*sqlcgen.BtcAccountKey) error`
- `UpdateMultisigAddr(accountType domainAccount.AccountType, item *sqlcgen.BtcAccountKey) (int64, error)`
- `UpdateMultisigAddrs(accountType domainAccount.AccountType, items []*sqlcgen.BtcAccountKey) (int64, error)`

**sqlcgen types used:** `sqlcgen.BtcAccountKey`

**Special features:**
- Transaction support in `UpdateMultisigAddrs` method
- Helper function `GetRedeemScriptByAddress(accountKeys []*sqlcgen.BtcAccountKey, addr string) string` in `interfaces.go`

**sqlcgen.BtcAccountKey structure:**
```go
- ID                  int64
- Coin                BtcAccountKeyCoin
- KeyType             string
- Account             BtcAccountKeyAccount
- P2pkhAddress        string
- P2shSegwitAddress   string
- Bech32Address       string
- TaprootAddress      sql.NullString
- FullPublicKey       string
- MultisigAddress     string
- RedeemScript        string
- WalletImportFormat  string
- Idx                 int64
- AddrStatus          int8
- UpdatedAt           sql.NullTime
```

**Security Note:** WalletImportFormat contains private key information

---

#### ETH Account Key Repository
**File:** `internal/infrastructure/repository/cold/eth_account_key_sqlc.go`
**Struct:** `ETHAccountKeyRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOneMaxID(accountType domainAccount.AccountType) (*sqlcgen.EthAccountKey, error)`
- `GetAllAddrStatus(accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus) ([]*sqlcgen.EthAccountKey, error)`
- `GetByAddress(address string) (*sqlcgen.EthAccountKey, error)`
- `InsertBulk(items []*sqlcgen.EthAccountKey) error`

**sqlcgen types used:** `sqlcgen.EthAccountKey`

**sqlcgen.EthAccountKey structure:**
```go
- ID              int64
- Account         EthAccountKeyAccount
- Address         string
- FullPublicKey   string
- PrivateKey      string
- Idx             int64
- AddrStatus      int8
- UpdatedAt       sql.NullTime
```

**Security Note:** PrivateKey contains sensitive cryptographic material

---

#### XRP Account Key Repository
**File:** `internal/infrastructure/repository/cold/xrp_account_key_sqlc.go`
**Struct:** `XRPAccountKeyRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetAllAddrStatus(ctx context.Context, accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus) ([]*sqlcgen.XrpAccountKey, error)`
- `InsertBulk(ctx context.Context, items []*sqlcgen.XrpAccountKey) error`

**sqlcgen types used:** `sqlcgen.XrpAccountKey`

**sqlcgen.XrpAccountKey structure:**
```go
- ID                int64
- Coin              XrpAccountKeyCoin
- Account           XrpAccountKeyAccount
- AccountID         string
- KeyType           int8
- MasterKey         string
- MasterSeed        string
- MasterSeedHex     string
- PublicKey         string
- PublicKeyHex      string
- IsRegularKeyPair  bool
- AllocatedID       sql.NullString
- AddrStatus        int8
- UpdatedAt         sql.NullTime
```

**Security Note:** MasterKey, MasterSeed, MasterSeedHex contain sensitive cryptographic material

---

#### Auth Account Key Repository
**File:** `internal/infrastructure/repository/cold/auth_account_key_sqlc.go`
**Struct:** `AuthAccountKeyRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(authType domainAccount.AuthType) (*sqlcgen.AuthAccountKey, error)`
- `Insert(item *sqlcgen.AuthAccountKey) error`

**sqlcgen types used:** `sqlcgen.AuthAccountKey`

**sqlcgen.AuthAccountKey structure:**
```go
- ID                  int16
- Coin                AuthAccountKeyCoin
- KeyType             string
- AuthAccount         string
- P2pkhAddress        string
- P2shSegwitAddress   string
- Bech32Address       string
- TaprootAddress      sql.NullString
- FullPublicKey       string
- MultisigAddress     string
- RedeemScript        string
- WalletImportFormat  string
- Idx                 int64
- AddrStatus          int8
- UpdatedAt           sql.NullTime
```

**Security Note:** WalletImportFormat contains private key information

---

#### Auth Full Pubkey Repository
**File:** `internal/infrastructure/repository/cold/auth_full_pubkey_sqlc.go`
**Struct:** `AuthFullPubkeyRepositorySqlc`

**Methods returning `sqlcgen` types:**
- `GetOne(authType domainAccount.AuthType) (*sqlcgen.AuthFullpubkey, error)`
- `InsertBulk(items []*sqlcgen.AuthFullpubkey) error`

**sqlcgen types used:** `sqlcgen.AuthFullpubkey`

**sqlcgen.AuthFullpubkey structure:**
```go
- ID              int16
- Coin            AuthFullpubkeyCoin
- AuthAccount     string
- FullPublicKey   string
- Fingerprint     sql.NullString
- UpdatedAt       sql.NullTime
```

---

#### Nonce Repository (EXEMPLARY - Already Refactored)
**File:** `internal/infrastructure/repository/cold/nonce_repository_sqlc.go`
**Struct:** `NonceRepositorySqlc`

**Methods returning domain types:**
- `GetNonce(ctx context.Context, signerID, transactionID string) (*multisig.NonceCommitment, error)`
- `GetAllNoncesForTransaction(ctx context.Context, transactionID string) ([]*multisig.NonceCommitment, error)`
- `GetUnusedNoncesForTransaction(ctx context.Context, transactionID string) ([]*multisig.NonceCommitment, error)`

**Pattern used:**
This repository already converts `sqlcgen.Musig2Nonce` to `multisig.NonceCommitment` domain entity internally using converter functions. **This is the pattern all other repositories should follow.**

**Example conversion pattern:**
```go
// Internal converter - sqlcgen to domain
func convertToNonceCommitment(nonce *sqlcgen.Musig2Nonce) (*multisig.NonceCommitment, error) {
    // Convert sqlcgen type to domain entity
}

// Public method returns domain entity
func (r *NonceRepositorySqlc) GetNonce(ctx context.Context, signerID, transactionID string) (*multisig.NonceCommitment, error) {
    nonce, err := r.queries.GetNonceBySignerAndTx(ctx, ...)
    return convertToNonceCommitment(&nonce)
}
```

---

## 2. Interface Migration Status

### 2.1 Current Status: Already Migrated

**All repository interfaces have been migrated** to `internal/application/ports/persistence/repository.go`

**Backward compatibility type aliases remain in:**
- `internal/infrastructure/repository/watch/interfaces.go`
- `internal/infrastructure/repository/cold/interfaces.go`

**Example from watch/interfaces.go:**
```go
// Type aliases for backward compatibility.
// These interfaces have been moved to pkg/application/ports/persistence.
type AddressRepositorier = persistence.AddressRepositorier
type BTCTxRepositorier = persistence.BTCTxRepositorier
type TxInputRepositorier = persistence.TxInputRepositorier
// ... etc
```

### 2.2 Action Required

**After refactoring is complete:**
- Remove type aliases from `watch/interfaces.go`
- Remove type aliases from `cold/interfaces.go`
- Migrate helper function `GetRedeemScriptByAddress` from cold/interfaces.go to domain layer

---

## 3. Existing Domain Entities

### 3.1 Domain Entity Locations

**Base Path:** `internal/domain/`

**Existing Domain Packages:**
- `account/` - Account and authentication types
- `address/` - Address types and status
- `bitcoin/` - Bitcoin-specific types
- `coin/` - Cryptocurrency types
- `key/` - Key generation types
- `multisig/` - Multisig and MuSig2 types (including `NonceCommitment`)
- `transaction/` - Transaction types
- `wallet/` - Wallet types

### 3.2 Key Domain Types (Already Defined)

**From `internal/domain/account/types.go`:**
- `AccountType` - Account purpose (client, deposit, payment, stored, auth1-15, anonymous, test)
- `AuthType` - Authorization types for multisig (auth1-auth15)

**From `internal/domain/address/types.go`:**
- `AddrType` - Address type (legacy, p2sh-segwit, bech32, taproot, bch-cashaddr, eth-address)
- `AddrStatus` - Address generation progress status

**From `internal/domain/transaction/types.go`:**
- `TxType` - Transaction lifecycle state (unsigned, signed, sent, done, notified, canceled)
- `ActionType` - Transaction operation type (deposit, payment, transfer)

**From `internal/domain/key/types.go`:**
- `KeyType` - Key generation standard (bip44, bip49, bip84, bip86, musig2)

**From `internal/domain/multisig/nonce_repository.go`:**
- `NonceCommitment` - MuSig2 nonce commitment entity

### 3.3 Domain Entity Gap Analysis

**Missing domain entities that need to be created:**

#### Watch Wallet Entities
1. **Address entity** (from `sqlcgen.Address`)
2. **Transaction entity** (from `sqlcgen.Tx`)
3. **BTCTransaction entity** (from `sqlcgen.BtcTx`)
4. **BTCTxInput entity** (from `sqlcgen.BtcTxInput`)
5. **BTCTxOutput entity** (from `sqlcgen.BtcTxOutput`)
6. **PaymentRequest entity** (from `sqlcgen.PaymentRequest`)
7. **ETHDetailTx entity** (from `sqlcgen.EthDetailTx`)
8. **XRPDetailTx entity** (from `sqlcgen.XrpDetailTx`)

#### Cold Wallet Entities
9. **Seed entity** (from `sqlcgen.Seed`)
10. **BTCAccountKey entity** (from `sqlcgen.BtcAccountKey`)
11. **ETHAccountKey entity** (from `sqlcgen.EthAccountKey`)
12. **XRPAccountKey entity** (from `sqlcgen.XrpAccountKey`)
13. **AuthAccountKey entity** (from `sqlcgen.AuthAccountKey`)
14. **AuthFullPubkey entity** (from `sqlcgen.AuthFullpubkey`)

---

## 4. Dependencies Between Components

### 4.1 Repository Dependencies

```
BTCTransaction Repository
    ├─ depends on: BTCTxInput Repository (for inputs)
    └─ depends on: BTCTxOutput Repository (for outputs)

PaymentRequest Repository
    └─ depends on: Address Repository (for address allocation)

Address Repository
    └─ independent (no repository dependencies)

Transaction Repository (tx_sqlc.go)
    └─ base transaction table, linked to coin-specific tx tables
```

### 4.2 Domain Entity Dependencies

```
Domain Entities (existing)
    ├─ AccountType
    ├─ AddrStatus
    ├─ TxType
    ├─ ActionType
    └─ KeyType

Domain Entities (to be created)
    ├─ Address
    │   └─ uses: AccountType
    │
    ├─ Transaction
    │   └─ uses: ActionType
    │
    ├─ BTCTransaction
    │   ├─ uses: TxType, ActionType
    │   ├─ related to: BTCTxInput
    │   └─ related to: BTCTxOutput
    │
    ├─ PaymentRequest
    │   └─ uses: AccountType
    │
    └─ AccountKey entities
        └─ uses: AccountType, AddrStatus, KeyType
```

---

## 5. Migration Plan with Sequencing

### 5.1 Phase 1: Foundation (Issue #225)
**Status:** This document

**Deliverables:**
- ✅ Comprehensive list of repository files returning `sqlcgen` types
- ✅ List of interfaces to migrate (already migrated, type aliases remain)
- ✅ Domain entity gap analysis
- ✅ Migration plan document with sequencing
- ✅ Risk assessment

---

### 5.2 Phase 2: Address and Transaction Core (Issue #226)
**Priority:** HIGH
**Estimated Size:** MEDIUM

**Scope:**
- `internal/infrastructure/repository/watch/address_sqlc.go`
- `internal/infrastructure/repository/watch/tx_sqlc.go`

**Domain entities to create:**
- `internal/domain/address/entity.go` - Address entity
- `internal/domain/transaction/entity.go` - Transaction entity

**Tasks:**
1. Create Address domain entity
2. Create Transaction domain entity
3. Add converter functions in address repository
4. Add converter functions in transaction repository
5. Update repository methods to return domain entities
6. Update application layer use cases
7. Write/update tests

**Dependencies:**
- None (foundational entities)

**Risk Level:** LOW
- Core entities with simple structure
- No complex relationships
- Establishes pattern for other repositories

---

### 5.3 Phase 3: BTC Transaction Details (Issue #227)
**Priority:** HIGH
**Estimated Size:** MEDIUM

**Dependencies:**
- Requires Phase 2 completion (Address and Transaction entities)
- Should follow pattern from Issue #226

**Scope:**
- `internal/infrastructure/repository/watch/btc_tx_sqlc.go`
- `internal/infrastructure/repository/watch/btc_tx_input_sqlc.go`
- `internal/infrastructure/repository/watch/btc_tx_output_sqlc.go`

**Domain entities to create:**
- `internal/domain/bitcoin/transaction.go` - BTCTransaction entity
- `internal/domain/bitcoin/tx_input.go` - BTCTxInput entity
- `internal/domain/bitcoin/tx_output.go` - BTCTxOutput entity

**Tasks:**
1. Create BTC transaction domain entities
2. Add converter functions in BTC repositories
3. Update repository methods to return domain entities
4. Handle UTXO model relationships (inputs/outputs)
5. Update BTC use cases
6. Write/update tests

**Dependencies:**
- Address entity (for input/output addresses)
- Transaction entity (base transaction)

**Risk Level:** MEDIUM
- Complex transaction structure
- UTXO model relationships
- Multi-sig transaction details

---

### 5.4 Phase 4: Payment Request (Future)
**Priority:** MEDIUM
**Estimated Size:** SMALL

**Scope:**
- `internal/infrastructure/repository/watch/payment_request_sqlc.go`

**Domain entities to create:**
- `internal/domain/payment/request.go` - PaymentRequest entity

**Dependencies:**
- Address entity (for sender/receiver addresses)
- AccountType (existing)

**Risk Level:** LOW
- Simple entity structure
- Transaction support needed

---

### 5.5 Phase 5: Multi-Chain Transaction Details (Future)
**Priority:** MEDIUM
**Estimated Size:** MEDIUM

**Scope:**
- `internal/infrastructure/repository/watch/eth_detail_tx_sqlc.go`
- `internal/infrastructure/repository/watch/xrp_detail_tx_sqlc.go`

**Domain entities to create:**
- `internal/domain/ethereum/transaction.go` - ETHDetailTx entity
- `internal/domain/ripple/transaction.go` - XRPDetailTx entity

**Dependencies:**
- Transaction entity (base transaction)
- Pattern from BTC transaction refactoring

**Risk Level:** MEDIUM
- Chain-specific transaction details
- Different transaction models from BTC

---

### 5.6 Phase 6: Cold Wallet Key Management (Future)
**Priority:** LOW (security-sensitive)
**Estimated Size:** LARGE

**Scope:**
- `internal/infrastructure/repository/cold/seed_sqlc.go`
- `internal/infrastructure/repository/cold/account_key_sqlc.go`
- `internal/infrastructure/repository/cold/eth_account_key_sqlc.go`
- `internal/infrastructure/repository/cold/xrp_account_key_sqlc.go`
- `internal/infrastructure/repository/cold/auth_account_key_sqlc.go`
- `internal/infrastructure/repository/cold/auth_full_pubkey_sqlc.go`

**Domain entities to create:**
- `internal/domain/key/seed.go` - Seed entity
- `internal/domain/bitcoin/account_key.go` - BTCAccountKey entity
- `internal/domain/ethereum/account_key.go` - ETHAccountKey entity
- `internal/domain/ripple/account_key.go` - XRPAccountKey entity
- `internal/domain/account/auth_key.go` - AuthAccountKey entity
- `internal/domain/account/auth_pubkey.go` - AuthFullPubkey entity

**Dependencies:**
- KeyType (existing)
- AddrStatus (existing)
- AccountType (existing)

**Risk Level:** HIGH
- **Security-sensitive:** Private keys, seeds, WIF
- Offline wallet operations (keygen, sign)
- Multi-sig key coordination
- Requires thorough security review
- Must maintain backward compatibility

---

## 6. Common Conversion Patterns

### 6.1 Null Type Conversions

**SQL null types to Go types:**
```go
// sql.NullTime → *time.Time or time.Time
if sqlcType.UpdatedAt.Valid {
    entity.UpdatedAt = &sqlcType.UpdatedAt.Time
}

// sql.NullString → *string or ""
if sqlcType.TaprootAddress.Valid {
    entity.TaprootAddress = &sqlcType.TaprootAddress.String
}

// sql.NullInt64 → *int64 or 0
if sqlcType.PaymentID.Valid {
    entity.PaymentID = &sqlcType.PaymentID.Int64
}
```

### 6.2 String Amount to BigInt

**For BTC amounts (string → *big.Int):**
```go
func convertAmount(amountStr string) (*big.Int, error) {
    if amountStr == "" {
        return big.NewInt(0), nil
    }
    amount := new(big.Int)
    amount, ok := amount.SetString(amountStr, 10)
    if !ok {
        return nil, fmt.Errorf("invalid amount: %s", amountStr)
    }
    return amount, nil
}
```

### 6.3 Enum Type Conversions

**Database int8 to domain enum:**
```go
// CurrentTxType (int8) → domain.TxType
func convertTxType(dbType int8) transaction.TxType {
    return transaction.TxType(dbType)
}

// AddrStatus (int8) → domain.AddrStatus
func convertAddrStatus(dbStatus int8) address.AddrStatus {
    return address.AddrStatus(dbStatus)
}
```

### 6.4 Coin/Account String to Domain Type

**String fields to domain types:**
```go
// Coin string → domain.CoinTypeCode
func convertCoin(coinStr string) coin.CoinTypeCode {
    return coin.CoinTypeCode(coinStr)
}

// Account string → domain.AccountType
func convertAccount(accountStr string) account.AccountType {
    return account.AccountType(accountStr)
}
```

---

## 7. Risk Assessment

### 7.1 Low Risk Components

**Characteristics:**
- Simple entity structure
- No complex relationships
- No security-sensitive data
- Easy to test

**Examples:**
- Address entity
- Transaction entity (base)
- PaymentRequest entity
- AuthFullPubkey entity

### 7.2 Medium Risk Components

**Characteristics:**
- Complex entity relationships
- UTXO or account model handling
- Chain-specific transaction logic

**Examples:**
- BTCTransaction with inputs/outputs
- ETHDetailTx entity
- XRPDetailTx entity

### 7.3 High Risk Components

**Characteristics:**
- **Security-sensitive data** (private keys, seeds)
- Offline wallet operations
- Multi-sig coordination
- Breaking changes to signing logic

**Examples:**
- Seed entity
- All AccountKey entities (BTC, ETH, XRP, Auth)
- Any repository handling WalletImportFormat or PrivateKey

**Required safeguards for high-risk components:**
1. Security review before implementation
2. Comprehensive test coverage
3. Backward compatibility verification
4. Offline wallet operation testing
5. Private key handling audit
6. No logging of sensitive fields

---

## 8. Design Principles Reference

**From `internal/AGENTS.md` (lines 245-380):**

### Core Principles:
1. **Repositories convert at the boundary:**
   - Accept domain entities as input
   - Return domain entities as output
   - Convert to/from database types internally

2. **No infrastructure types leak to application layer:**
   - Use cases work with domain entities only
   - No `sqlcgen` types in application layer

3. **Domain layer has zero infrastructure dependencies:**
   - Domain entities are pure Go structs
   - No database annotations
   - No external package dependencies

4. **Converter functions are private:**
   - Keep conversion logic in repository implementation
   - Only expose domain entity methods publicly

---

## 9. Reference Implementation

**File:** `internal/infrastructure/repository/cold/nonce_repository_sqlc.go`

**Why this is exemplary:**
- ✅ Returns domain entity (`multisig.NonceCommitment`)
- ✅ Private converter functions
- ✅ Clean separation of concerns
- ✅ Proper error handling in conversions

**Key code pattern:**
```go
// Private converter
func convertToNonceCommitment(nonce *sqlcgen.Musig2Nonce) (*multisig.NonceCommitment, error) {
    pubNonceBytes, err := hex.DecodeString(nonce.PubNonce)
    if err != nil {
        return nil, fmt.Errorf("failed to decode pub_nonce: %w", err)
    }

    secNonceBytes, err := hex.DecodeString(nonce.SecNonce)
    if err != nil {
        return nil, fmt.Errorf("failed to decode sec_nonce: %w", err)
    }

    return &multisig.NonceCommitment{
        SignerID:      nonce.SignerID,
        TransactionID: nonce.TransactionID,
        PubNonce:      pubNonceBytes,
        SecNonce:      secNonceBytes,
        Used:          nonce.Used,
        CreatedAt:     nonce.CreatedAt,
    }, nil
}

// Public method returns domain entity
func (r *NonceRepositorySqlc) GetNonce(ctx context.Context, signerID, transactionID string) (*multisig.NonceCommitment, error) {
    nonce, err := r.queries.GetNonceBySignerAndTx(ctx, sqlcgen.GetNonceBySignerAndTxParams{
        SignerID:      signerID,
        TransactionID: transactionID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get nonce: %w", err)
    }

    return convertToNonceCommitment(&nonce)
}
```

---

## 10. Critical Notes

### 10.1 Security Considerations

**Private key fields requiring special handling:**
- `WalletImportFormat` (BTC, Auth account keys)
- `PrivateKey` (ETH account keys)
- `MasterKey`, `MasterSeed`, `MasterSeedHex` (XRP account keys)
- `Seed` (Seed repository)

**Security requirements:**
- Never log these fields
- Consider encryption at rest
- Audit all code paths handling these fields
- Security review required before merge

### 10.2 Transaction Support

**Repositories with database transaction support:**
- `TxRepositorySqlc` - implements `WithTx()`
- `PaymentRequestRepositorySqlc` - implements `WithTx()`
- `ETHDetailTXRepositorySqlc` - implements `WithTx()`
- `XRPDetailTXRepositorySqlc` - implements `WithTx()`
- `BTCAccountKeyRepositorySqlc` - transaction in `UpdateMultisigAddrs()`

**Action required:**
- Maintain transaction support after refactoring
- Test transactional operations thoroughly

### 10.3 Helper Functions to Migrate

**From `internal/infrastructure/repository/cold/interfaces.go`:**
```go
func GetRedeemScriptByAddress(accountKeys []*sqlcgen.BtcAccountKey, addr string) string
```

**Action required:**
- Migrate to domain layer or application service
- Update to use domain entity instead of `sqlcgen.BtcAccountKey`

### 10.4 Backward Compatibility

**Type aliases to remove after refactoring:**
- `internal/infrastructure/repository/watch/interfaces.go`
- `internal/infrastructure/repository/cold/interfaces.go`

**Timing:**
- Remove after all repositories are refactored
- Ensure no callers use these aliases
- Update imports across codebase

---

## 11. Success Criteria

### Completion Criteria for Each Phase:

1. **Domain entities created** with proper field types
2. **Converter functions implemented** and tested
3. **Repository methods updated** to return domain entities
4. **Interface definitions updated** in `application/ports/persistence/`
5. **All use cases updated** to work with domain entities
6. **All tests passing:**
   - `make lint-fix` passes
   - `make tidy` passes
   - `make check-build` passes
   - `make gotest` passes
7. **No breaking changes** to existing functionality
8. **Security review completed** (for high-risk components)
9. **Documentation updated** (godoc comments, architecture diagrams)

---

## 12. Conclusion

This audit provides a complete foundation for refactoring the infrastructure layer:

**Documented:**
- ✅ 19 repository files requiring refactoring
- ✅ 14 distinct sqlcgen types
- ✅ Interface migration status
- ✅ Existing domain entities
- ✅ 14 missing domain entities identified
- ✅ Clear migration plan with 6 phases
- ✅ Risk assessment for each component
- ✅ Reference implementation pattern
- ✅ Security considerations
- ✅ Common conversion patterns

**Next Steps:**
1. Review and approve this audit document
2. Create domain entities for Address and Transaction (Issue #226)
3. Refactor address and transaction repositories (Issue #226)
4. Create domain entities for BTC transactions (Issue #227)
5. Refactor BTC transaction repositories (Issue #227)
6. Continue with remaining phases as prioritized

**Key Takeaway:**
The `NonceRepositorySqlc` implementation provides an excellent pattern to follow. By applying this pattern systematically across all repositories, the codebase will achieve true Clean Architecture separation with domain entities properly isolated from infrastructure concerns.

---

**Document Version:** 1.0
**Last Updated:** 2026-01-02
**Author:** Infrastructure Refactoring Team
**Review Status:** Pending
