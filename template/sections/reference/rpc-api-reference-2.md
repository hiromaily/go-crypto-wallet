## RPC & API Reference

### Bitcoin Cash Node RPC

**Essential Commands:**

| Command | Description |
|---------|-------------|
| `getblockchaininfo` | Network and sync status |
| `getbalance` | Wallet balance |
| `listunspent` | List UTXOs |
| `createrawtransaction` | Create raw transaction |
| `signrawtransactionwithkey` | Sign with provided keys |
| `sendrawtransaction` | Broadcast transaction |
| `gettransaction` | Get transaction details |
| `getaddressinfo` | Get address information |
| `validateaddress` | Validate address format |

### BCH-Specific RPC Differences

| Command | BCH Behavior | Notes |
|---------|--------------|-------|
| `getaddressinfo` | Different response structure | No SegWit fields |
| `decodescript` | No witness fields | Legacy scripts only |

### Go Libraries

| Library | Purpose | Repository |
|---------|---------|------------|
| **gcash/bchd** | BCH full node in Go | [github.com/gcash/bchd](https://github.com/gcash/bchd) |
| **gcash/bchutil** | BCH utilities | [github.com/gcash/bchutil](https://github.com/gcash/bchutil) |
| **cpacia/bchutil** | CashAddr encoding | [github.com/cpacia/bchutil](https://github.com/cpacia/bchutil) |

### CashAddr Library Usage

```go
import "github.com/cpacia/bchutil"

// Encode to CashAddr
cashaddr, err := bchutil.EncodeAddress(hash160, &chaincfg.MainNetParams)

// Decode from CashAddr
decoded, err := bchutil.DecodeAddress(cashaddr, &chaincfg.MainNetParams)
```

---
