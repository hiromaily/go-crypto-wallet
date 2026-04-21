## Network & Consensus

### Networks

| Network | Purpose | P2P Port | RPC Port |
|---------|---------|----------|----------|
| **Mainnet** | Production | 8333 | 8332 |
| **Testnet3** | Public testing | 18333 | 18332 |
| **Testnet4** | Public testing (newer) | 28333 | 28332 |
| **Regtest** | Local development | 18444 | 18443 |
| **Chipnet** | Contract testing | 48333 | 48332 |

### Block Structure

```
+----------------------+
| Block Header (80 B)  |
|   - Version (4 B)    |
|   - PrevBlockHash    |  32 bytes
|   - MerkleRoot       |  32 bytes
|   - Timestamp (4 B)  |
|   - Bits (4 B)       |  Difficulty target
|   - Nonce (4 B)      |
+----------------------+
| Transaction Count    |  VarInt
+----------------------+
| Transactions[]       |  Up to 32 MB
+----------------------+
```

### Difficulty Adjustment Algorithm (DAA)

BCH uses ASERT (Absolutely Scheduled Exponentially Rising Targets) DAA since November 2020:

- Adjusts difficulty every block
- Targets 10-minute block times
- Responds smoothly to hashrate changes

### Node Implementations

| Implementation | Description | Repository |
|----------------|-------------|------------|
| **Bitcoin Cash Node (BCHN)** | Primary implementation | [bitcoincashnode/bitcoincashnode](https://gitlab.com/bitcoin-cash-node/bitcoin-cash-node) |
| **Bitcoin Unlimited** | Alternative implementation | [BitcoinUnlimited/BitcoinUnlimited](https://github.com/BitcoinUnlimited/BitcoinUnlimited) |
| **BCHD** | Go implementation | [gcash/bchd](https://github.com/gcash/bchd) |
| **Flowee** | C++ implementation | [flowee-org/thehub](https://gitlab.com/FloweeTheHub/thehub) |

### Recommended Node Version

- **Bitcoin Cash Node**: v27.0.0+ (2026)
- Includes latest consensus rules and CashTokens support

---
