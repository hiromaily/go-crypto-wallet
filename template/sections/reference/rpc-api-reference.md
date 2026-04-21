## RPC & API Reference

### Bitcoin Core RPC

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
| `walletprocesspsbt` | Process PSBT |
| `finalizepsbt` | Finalize PSBT |
| `decodepsbt` | Decode/analyze PSBT |

**Reference:**

- [Bitcoin Core RPC Documentation](https://developer.bitcoin.org/reference/rpc/)
- [Bitcoin Core JSON-RPC API Reference](https://bitcoincore.org/en/doc/)

### Go Libraries

| Library | Purpose | Repository |
|---------|---------|------------|
| **btcd** | Full node implementation | [github.com/btcsuite/btcd](https://github.com/btcsuite/btcd) |
| **btcutil** | Address/transaction utilities | [github.com/btcsuite/btcd/btcutil](https://github.com/btcsuite/btcd) |
| **btcec** | secp256k1 cryptography | [github.com/btcsuite/btcd/btcec](https://github.com/btcsuite/btcd) |
| **psbt** | PSBT implementation | [github.com/btcsuite/btcd/btcutil/psbt](https://github.com/btcsuite/btcd) |
| **txscript** | Script parsing/building | [github.com/btcsuite/btcd/txscript](https://github.com/btcsuite/btcd) |

---
