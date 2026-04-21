## RPC & API Reference

### Key go-ethereum Methods Used

| Method | Description |
|--------|-------------|
| `ethclient.BalanceAt` | Get account balance at block |
| `ethclient.PendingNonceAt` | Get pending nonce |
| `ethclient.SendTransaction` | Broadcast raw transaction |
| `eth_gasPrice` | Get current gas price (legacy) |
| `eth_estimateGas` | Estimate gas for transaction |
| `eth_getTransactionCount` | Get nonce (use `pending` tag) |
| `eth_getTransactionReceipt` | Get receipt (for confirmation status) |
| `eth_getBlockByNumber` | Get block (for baseFeePerGas detection) |
| `eth_call` | Call contract function (read-only) |
| `net_version` | Get network/chain ID |

### Go Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| **go-ethereum** | v1.17.0 | Ethereum client library |
| `ethclient` | — | HTTP/WS client |
| `core/types` | — | Transaction types (LegacyTx, DynamicFeeTx) |
| `crypto` | — | Keccak256, PubkeyToAddress |
| `accounts/keystore` | — | Encrypted key storage |
| `common` | — | HexToAddress, IsHexAddress |
| `rlp` | — | RLP encode/decode |
| **btcd/btcec/v2** | — | ECDSA private key operations |
| **btcd/btcutil/hdkeychain** | — | BIP44 HD key derivation |

### Application Port Interface

The `Ethereumer` interface (defined in `internal/application/ports/api/eth/interface.go`) provides:

- Balance operations: `GetTotalBalance`, `BalanceAt`
- Network: `Syncing`, `ProtocolVersion`, `NetVersion`
- Block: `BlockNumber`, `GetBlockByNumber`, `EnsureBlockNumber`
- Transaction: `CreateRawTransaction`, `CreateRawTransactionEIP1559`, `SendRawTransaction`
- Gas: `GasPrice`, `EstimateGas`
- Signing: `SignOnRawTransaction`
- Key management: `ToECDSA`, `GetPrivKey`, `ImportRawKey`

---
