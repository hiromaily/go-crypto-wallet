## ETH-Specific Flow Details

> The common 3-wallet setup, signing, and monitoring flows are defined in the chain-agnostic reference:
> [docs/transaction-flow.md](../../../docs/transaction-flow.md).
> This section describes Ethereum-specific concerns on top of that common flow.

### Single-Sig Flow (Ethereum)

Follows the [common single-sig flow](../../../docs/transaction-flow.md#single-sig-flow).

Ethereum-specific steps:

- Watch Wallet fetches the pending nonce via `eth_getTransactionCount` with `pending` tag
- Watch Wallet fetches gas price via `eth_gasPrice` and estimates gas via `eth_estimateGas`
- Watch Wallet serializes the transaction as RLP-encoded hex
- Keygen Wallet decodes the hex, signs with `types.NewLondonSigner(chainID)`, re-encodes to hex
- Watch Wallet sends via `eth_sendRawTransaction`
- Watch Wallet polls for receipt via `eth_getTransactionReceipt` to confirm finality

### Transaction Types vs Use Cases

| Use Case | Watch Action | Nonce Source | Amount Calculation |
|----------|-------------|-------------|-------------------|
| **Deposit** | Client address → deposit address | Pending nonce of client address | `balance - fee` |
| **Payment** | Payment address → external address | Pending nonce of payment address | Fixed amount from payment_request |
| **Transfer** | deposit/stored → other internal | Pending nonce of source address | `balance - fee` or fixed amount |

### Multisig

Ethereum EOA does not natively support multisig at the protocol level (unlike Bitcoin's P2SH/P2WSH). This system implements multisig via **Gnosis Safe v1.4.1** — an audited smart contract that enforces an m-of-n threshold before executing any transaction.

The implementation uses a file-based, offline-signing workflow:

1. Watch Wallet proposes a transaction and writes an unsigned JSON file
2. Each owner (Keygen or Sign wallet) verifies the EIP-712 `safeTxHash` offline and appends a signature
3. When the threshold is reached, Watch Wallet submits `execTransaction` on-chain

See [multisig.md](../../../docs/chains/eth/multisig.md) for the complete reference including file format, CLI commands, EIP-712 signing details, and E2E Pattern 3 (2-of-2 Safe payment).

---

**Document Version:** 1.1
**Last Updated:** 2026-03-08
**Maintainer:** go-crypto-wallet team
