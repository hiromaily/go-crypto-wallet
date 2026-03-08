# Transaction Flow

This document describes the common transaction flow that applies across all supported chains (BTC, BCH, ETH, XRP).
For chain-specific details, see the documentation under `docs/chains/`.

## Architecture Overview

The system uses three wallet types running on **separate machines**:

| Wallet | Environment | Role |
|--------|-------------|------|
| **Watch Wallet** | Online | Creates unsigned transactions, broadcasts signed transactions, monitors status |
| **Keygen Wallet** | Offline (air-gapped) | Generates keys and addresses; **sole signer for single-sig transactions** |
| **Sign Wallet** | Offline (air-gapped) | Provides additional signatures for multisig transactions only |

> **Single-sig vs Multisig wallet requirements:**
> - **Single-sig**: Only **Watch + Keygen** wallets are needed. The Sign Wallet is not required.
> - **Multisig (M-of-N)**: All three wallets are used. Keygen provides the first signature; Sign Wallet(s) provide the remaining signatures.

Key design principles:

- Only Watch Wallet connects to the network — Keygen and Sign Wallets are always offline
- Transaction files are transferred between machines via secure offline methods (e.g., USB drives)
- Each wallet maintains its own independent database

---

## 1. Setup Flow

Before sending any transaction, each wallet must be initialized and addresses must be shared with Watch Wallet.

### Single-sig Setup

```mermaid
sequenceDiagram
    participant User
    participant Keygen as Keygen Wallet<br/>(Offline)
    participant Watch as Watch Wallet<br/>(Online)
    participant File as File Transfer<br/>(USB, etc.)

    Note over User,File: Step 1: Keygen Wallet Initialization
    User->>Keygen: create seed
    User->>Keygen: create hdkey --account <account>

    Note over User,File: Step 2: Export and Import Addresses
    User->>Keygen: export address --account <account>
    Keygen->>File: Write address file
    User->>File: Transfer address file to Watch Wallet machine
    User->>Watch: import address --file <address_file>
    Watch->>Watch: Start monitoring imported addresses
```

### Multisig Setup (BTC / BCH / XRP)

> **ETH Note:** Ethereum Safe multisig does not use BIP32 multisig address derivation or fullpubkey exchange.
> Instead, a Safe smart contract is deployed on-chain with a list of owner EOA addresses.
> See [ETH Safe Multisig Setup](#eth-safe-multisig-setup) below.

```mermaid
sequenceDiagram
    participant User
    participant Keygen as Keygen Wallet<br/>(Offline)
    participant Sign1 as Sign Wallet 1<br/>(Offline)
    participant Sign2 as Sign Wallet 2<br/>(Offline)
    participant Watch as Watch Wallet<br/>(Online)
    participant File as File Transfer<br/>(USB, etc.)

    Note over User,File: Step 1: Keygen Wallet Initialization
    User->>Keygen: create seed
    User->>Keygen: create hdkey --account <account>

    Note over User,File: Step 2: Sign Wallet Initialization (each signer)
    User->>Sign1: create seed
    User->>Sign1: create hdkey
    User->>Sign1: export fullpubkey
    Sign1->>File: Write fullpubkey file

    User->>Sign2: create seed
    User->>Sign2: create hdkey
    User->>Sign2: export fullpubkey
    Sign2->>File: Write fullpubkey file

    Note over User,File: Step 3: Create Multisig Address in Keygen
    User->>File: Transfer fullpubkey files to Keygen machine
    User->>Keygen: import fullpubkey --file <sign1_fullpubkey>
    User->>Keygen: import fullpubkey --file <sign2_fullpubkey>
    User->>Keygen: create multisig --account <account>

    Note over User,File: Step 4: Share Addresses with Watch Wallet
    User->>Keygen: export address --account <account>
    Keygen->>File: Write address file
    User->>File: Transfer address file to Watch Wallet machine
    User->>Watch: import address --file <address_file>
    Watch->>Watch: Start monitoring imported addresses
```

### ETH Safe Multisig Setup

Ethereum multisig is implemented via Gnosis Safe v1.4.1. Each owner is a regular EOA. Setup requires deploying a Safe proxy contract on-chain with the owner list and threshold.

```mermaid
sequenceDiagram
    participant User
    participant Keygen as Keygen Wallet<br/>(Offline)
    participant Sign1 as Sign Wallet 1<br/>(Offline)
    participant Watch as Watch Wallet<br/>(Online)
    participant Chain as Ethereum Network
    participant File as File Transfer<br/>(USB, etc.)

    Note over User,File: Step 1: Generate owner EOA keys (one per signer wallet)
    User->>Keygen: create seed / create hdkey / export address
    Keygen->>File: Write keygen address file
    User->>Sign1: create seed / create hdkey / export address
    Sign1->>File: Write sign address file

    Note over User,File: Step 2: Import addresses into Watch Wallet (for monitoring)
    User->>Watch: import address --file <keygen_address_file>
    User->>Watch: import address --file <sign_address_file>

    Note over User,Chain: Step 3: Deploy Safe contract on-chain
    User->>Chain: forge script DeploySafe.s.sol<br/>(owners=[keygen_addr, sign_addr], threshold=2)
    Chain-->>User: Safe proxy address

    Note over User,Watch: Step 4: Configure Safe address
    User->>Watch: watch safe info --safe-address <safe_addr>
```

---

## 2. Transaction Operation Flow

### Single-sig Flow

Used when the address is controlled by a single key (Keygen Wallet only).

```mermaid
sequenceDiagram
    participant User
    participant Watch as Watch Wallet<br/>(Online)
    participant Keygen as Keygen Wallet<br/>(Offline)
    participant Blockchain as Blockchain Network
    participant File as File Transfer<br/>(USB, etc.)

    Note over User,File: Step 1: Create Unsigned Transaction (Watch Wallet)
    User->>Watch: create deposit / payment / transfer
    Watch->>Watch: Build unsigned transaction
    Watch->>File: Write unsigned tx file
    Watch-->>User: Return: [unsigned tx], [fileName]

    Note over User,File: Step 2: Sign Transaction (Keygen Wallet)
    User->>File: Transfer unsigned tx file to Keygen machine
    User->>Keygen: sign --file <unsigned_tx_file>
    Keygen->>Keygen: Sign with private key
    Keygen->>File: Write signed tx file
    Keygen-->>User: Return: [signed tx], [isCompleted: true], [fileName]

    Note over User,File: Step 3: Broadcast Transaction (Watch Wallet)
    User->>File: Transfer signed tx file to Watch machine
    User->>Watch: send --file <signed_tx_file>
    Watch->>Blockchain: Broadcast transaction
    Blockchain-->>Watch: Return: txID
    Watch-->>User: Return: txID
```

### Multisig Flow (BTC / BCH / XRP — M-of-N)

Used when the address requires multiple signatures. Signing is repeated until the required threshold M is met.

> **ETH Note:** Ethereum Safe multisig uses a different signing flow.
> See [ETH Safe Multisig Flow](#eth-safe-multisig-flow) below.

```mermaid
sequenceDiagram
    participant User
    participant Watch as Watch Wallet<br/>(Online)
    participant Keygen as Keygen Wallet<br/>(Offline)
    participant Sign1 as Sign Wallet 1<br/>(Offline)
    participant SignN as Sign Wallet N<br/>(Offline)
    participant Blockchain as Blockchain Network
    participant File as File Transfer<br/>(USB, etc.)

    Note over User,File: Step 1: Create Unsigned Transaction (Watch Wallet)
    User->>Watch: create deposit / payment / transfer
    Watch->>Watch: Build unsigned transaction
    Watch->>File: Write unsigned tx file
    Watch-->>User: Return: [unsigned tx], [fileName]

    Note over User,File: Step 2: First Signature (Keygen Wallet)
    User->>File: Transfer unsigned tx file to Keygen machine
    User->>Keygen: sign --file <unsigned_tx_file>
    Keygen->>Keygen: Sign with private key (1st signature)
    Keygen->>File: Write partially signed tx file
    Keygen-->>User: Return: [partial tx], [isCompleted: false], [fileName]

    Note over User,File: Step 3+: Additional Signatures (Sign Wallets, repeat until isCompleted: true)
    User->>File: Transfer partially signed tx file to Sign Wallet 1 machine
    User->>Sign1: sign --file <partial_tx_file>
    Sign1->>Sign1: Sign with private key (2nd signature)
    Sign1->>File: Write partially signed tx file
    Sign1-->>User: Return: [partial tx], [isCompleted: false or true], [fileName]

    Note over User,File: Continue signing with Sign Wallet N if more signatures are needed
    User->>File: Transfer partially signed tx file to Sign Wallet N machine
    User->>SignN: sign --file <partial_tx_file>
    SignN->>SignN: Sign with private key (Nth signature)
    SignN->>File: Write fully signed tx file
    SignN-->>User: Return: [full tx], [isCompleted: true], [fileName]

    Note over User,File: Final Step: Broadcast Transaction (Watch Wallet)
    User->>File: Transfer fully signed tx file to Watch machine
    User->>Watch: send --file <fully_signed_tx_file>
    Watch->>Blockchain: Broadcast transaction
    Blockchain-->>Watch: Return: txID
    Watch-->>User: Return: txID
```

> **Note**: The signing loop continues until `isCompleted: true` is returned. For a 2-of-3 multisig,
> only 2 signatures are needed; Sign Wallet 2 can be skipped once the threshold is met.

### ETH Safe Multisig Flow

Ethereum multisig uses Gnosis Safe v1.4.1. Key differences from BTC-style multisig:

- **Proposal**: Watch Wallet calls `watch create multisig` (not `create deposit/payment/transfer`) to produce an unsigned JSON proposal file containing an EIP-712 `safeTxHash`
- **Signing**: Each signer calls `sign signature --signer-address` (keygen or sign wallet). The wallet recomputes the `safeTxHash` offline from the file fields and appends a 65-byte EIP-712 signature
- **Submission**: Watch Wallet calls `watch send multisig send-eth` which passes the packed signatures to `execTransaction` on the Safe contract — not `eth_sendRawTransaction`
- **File counter**: Each signing step writes a new file with an incremented counter (`_0.json` → `_1.json` → `_2.json`)

```mermaid
sequenceDiagram
    participant User
    participant Watch as Watch Wallet<br/>(Online)
    participant Keygen as Keygen Wallet<br/>(Offline)
    participant Sign1 as Sign Wallet 1<br/>(Offline)
    participant Chain as Ethereum Network
    participant File as File Transfer<br/>(USB, etc.)

    Note over User,File: Step 1: Propose Multisig Transaction (Watch Wallet)
    User->>Watch: watch create multisig --safe-address <safe> --to <addr> --amount <eth>
    Watch->>Chain: getTransactionHash() — fetch safeTxHash via eth_call
    Chain-->>Watch: safeTxHash
    Watch->>File: Write {action}_multisig_{uuid}_0.json
    Watch-->>User: Return: [filePath], [uuid]

    Note over User,File: Step 2: First Signature (Keygen Wallet)
    User->>File: Transfer _0.json to Keygen machine
    User->>Keygen: keygen sign signature --file <_0.json> --signer-address <addr>
    Keygen->>Keygen: Recompute safeTxHash offline (EIP-712), verify, sign
    Keygen->>File: Write {action}_multisig_{uuid}_1.json
    Keygen-->>User: Return: [filePath], [isComplete: false], [signCount: 1]

    Note over User,File: Step 3: Second Signature (Sign Wallet — threshold = 2)
    User->>File: Transfer _1.json to Sign Wallet machine
    User->>Sign1: sign sign signature --file <_1.json> --signer-address <addr>
    Sign1->>Sign1: Recompute safeTxHash offline (EIP-712), verify, sign
    Sign1->>File: Write {action}_multisig_{uuid}_2.json (tx_type: "signed")
    Sign1-->>User: Return: [filePath], [isComplete: true], [signCount: 2]

    Note over User,File: Step 4: Submit to Safe Contract (Watch Wallet)
    User->>File: Transfer _2.json to Watch machine
    User->>Watch: watch send multisig send-eth --file <_2.json>
    Watch->>Chain: execTransaction(to, value, data, ..., packedSignatures)
    Chain-->>Watch: tx receipt
    Watch-->>User: Return: txHash
```

For full details see [`docs/chains/eth/multisig.md`](chains/eth/multisig.md).

---

## 3. Transaction Types

All transaction types follow the same signing flow above. They differ only in how Watch Wallet selects UTXOs/sources and destinations.

| Type | Description | Typical Signing |
|------|-------------|-----------------|
| **Deposit** | Consolidates funds received at client addresses into a managed cold wallet address | Single-sig or Multisig |
| **Payment** | Sends funds to an external address based on a withdrawal request | Multisig (recommended) |
| **Transfer** | Moves funds between internal accounts (e.g., deposit → payment account) | Single-sig or Multisig |

---

## 4. Monitoring Flow

After broadcasting, Watch Wallet tracks confirmation status automatically.

```mermaid
sequenceDiagram
    participant User
    participant Watch as Watch Wallet<br/>(Online)
    participant Blockchain as Blockchain Network

    Note over User,Blockchain: Manual status check
    User->>Watch: monitor senttx --account <account>
    Watch-->>User: Display sent transaction list and status

    Note over User,Blockchain: Automatic confirmation tracking (periodic)
    loop Until sufficient confirmations
        Watch->>Blockchain: Query transaction status (txID)
        Blockchain-->>Watch: Return: confirmations

        alt Confirmations >= threshold
            Watch->>Watch: Update status: Sent → Done
            Watch->>Watch: Send notification
            Watch->>Watch: Update status: Done → Notified
        else Confirmations < threshold
            Watch->>Watch: Continue monitoring
        end
    end
```

### Transaction Status Lifecycle

```
TxTypeSent ──(confirmations >= threshold)──> TxTypeDone ──(notified)──> TxTypeNotified
```

---

## 5. Security Model

| Principle | Description |
|-----------|-------------|
| **Network isolation** | Keygen and Sign Wallets never connect to the internet |
| **Private key separation** | Watch Wallet holds no private keys; it only works with public addresses |
| **Multisig** | Multiple independent offline machines must cooperate to authorize transactions |
| **Offline file transfer** | Transaction files move between machines via physically controlled media |
| **Role separation** | Key generation, transaction signing, and broadcasting are handled by distinct wallets |

---

## Chain-Specific Documentation

For chain-specific details (address formats, transaction formats, signing algorithms, etc.), see:

| Chain | Documentation |
|-------|---------------|
| Bitcoin (BTC) | [`docs/chains/btc/operations/wallet-flow.md`](chains/btc/operations/wallet-flow.md) |
| Bitcoin Cash (BCH) | [`docs/chains/bch/README.md`](chains/bch/README.md) |
| Ethereum (ETH) | [`docs/chains/eth/README.md`](chains/eth/README.md) |
| Ripple (XRP) | [`docs/chains/xrp/README.md`](chains/xrp/README.md) |
