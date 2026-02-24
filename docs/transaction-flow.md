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

### Multisig Setup

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

### Multisig Flow (M-of-N)

Used when the address requires multiple signatures. Signing is repeated until the required threshold M is met.

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
