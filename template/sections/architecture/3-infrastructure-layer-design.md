## 3. Infrastructure Layer Design

### 3.1 Interface Definition

```go
package btc

import (
    "github.com/btcsuite/btcd/btcutil/psbt"
    "github.com/btcsuite/btcd/wire"
)

// PSBTOperator defines PSBT operations interface
type PSBTOperator interface {
    // Creation (Watch Wallet)
    CreatePSBTFromTx(msgTx *wire.MsgTx, prevTxs []PrevTx) (string, error)

    // Parsing
    ParsePSBT(psbtBase64 string) (*ParsedPSBT, error)

    // Validation
    ValidatePSBT(psbtBase64 string) error

    // Signing (offline) - all metadata in PSBT per BIP174
    SignPSBTWithKey(psbtBase64 string, wifs []string) (string, bool, error)

    // Finalization (Watch Wallet)
    FinalizePSBT(psbtBase64 string) error

    // Extraction (Watch Wallet)
    ExtractTransaction(psbtBase64 string) (*wire.MsgTx, error)

    // Utility
    IsPSBTComplete(psbtBase64 string) (bool, error)
    GetPSBTFee(psbtBase64 string) (int64, error)
}

// ParsedPSBT contains parsed PSBT information
type ParsedPSBT struct {
    UnsignedTx      *wire.MsgTx
    Inputs          []PSBTInput
    Outputs         []PSBTOutput
    IsComplete      bool
    IsPartiallySigned bool
    SignatureCount  int
    RequiredSigs    int
}

// PSBTInput represents per-input data
type PSBTInput struct {
    PrevTxID        string
    PrevVout        uint32
    PrevAmount      int64
    PrevScriptPubKey []byte
    RedeemScript    []byte
    WitnessScript   []byte
    Signatures      [][]byte
    PublicKeys      [][]byte
}

// PSBTOutput represents per-output data
type PSBTOutput struct {
    Address         string
    Amount          int64
    ScriptPubKey    []byte
    RedeemScript    []byte
    WitnessScript   []byte
}
```

### 3.2 Implementation Structure

```
internal/infrastructure/api/btc/btc/
├── psbt.go              (new) - PSBT operations implementation
├── psbt_rpc.go          (new) - Bitcoin Core RPC PSBT methods
├── psbt_offline.go      (new) - Offline PSBT signing (btcd)
├── psbt_test.go         (new) - Unit tests
├── transaction.go              - Existing transaction methods
└── bitcoin.go                  - Bitcoin client
```

### 3.3 Method Responsibilities

**psbt.go** - Main PSBT interface

- Interface definition
- Common PSBT utilities
- Validation functions

**psbt_rpc.go** - RPC methods (Watch Wallet)

- `CreatePSBTFromTx` using `walletcreatefundedpsbt`
- `FinalizePSBT` using `finalizepsbt`
- `CombinePSBT` using `combinepsbt`

**psbt_offline.go** - Offline methods (Keygen/Sign Wallets)

- `SignPSBTWithKey` using btcd package
- Signature creation
- Private key operations

---
