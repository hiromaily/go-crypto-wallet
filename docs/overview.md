# Overview

This repository is a wallet system that handles cryptocurrencies such as Bitcoin, Bitcoin Cash, Ethereum, ERC-20 Tokens, and Ripple.
With a security-focused design, it separates three different types of wallets (Watch Wallet, Keygen Wallet, and Sign Wallet), each serving a different role.

## Wallet Types

### 1. Watch Wallet

**Features:**

- Operates in an **online environment**
- Holds **public keys only** (does not hold private keys)
- Can access blockchain nodes

**Main Functions:**

- Create unsigned transactions
- Send signed transactions
- Monitor transaction status
- Import addresses (public key addresses exported from Keygen Wallet)

**Main CLI Commands:**

- `watch import address` - Import addresses
- `watch create deposit` - Create unsigned deposit transactions
- `watch create payment` - Create unsigned payment transactions
- `watch create transfer` - Create unsigned transfer transactions between accounts
- `watch send` - Send signed transactions
- `watch monitor` - Monitor transactions and balances
- `watch api` - API calls to blockchain nodes (BTC/BCH/ETH/XRP specific)

### 2. Keygen Wallet

**Features:**

- Operates in an **offline environment** (cold wallet)
- Manages account key management
- Key generation based on HD Wallet (Hierarchical Deterministic Wallet)

**Main Functions:**

- Generate seeds for accounts
- Generate keys based on HD Wallet
- Generate multisig addresses (based on account configuration)
- Export public key addresses (CSV format, for import to Watch Wallet)
- **First signature** on unsigned transactions (for multisig)

**Main CLI Commands:**

- `keygen create seed` - Generate seeds
- `keygen create hdkey` - Generate HD keys
- `keygen create multisig` - Generate multisig addresses
- `keygen export address` - Export public key addresses
- `keygen import full-pubkey` - Import full public keys exported from Sign Wallet
- `keygen sign` - Sign unsigned transactions (first signature)
- `keygen api` - Wallet management API (BTC/BCH/ETH specific)

### 3. Sign Wallet

**Features:**

- Operates in an **offline environment** (cold wallet)
- Used by authentication operators
- Each operator is provided with their own authentication account and Sign Wallet application

**Main Functions:**

- Generate seeds for authentication accounts
- Generate HD keys for authentication accounts
- Export full public key addresses (CSV format, for import to Keygen Wallet)
- **Second and subsequent signatures** on unsigned transactions (for multisig)

**Main CLI Commands:**

- `sign create seed` - Generate seeds
- `sign create hdkey` - Generate HD keys
- `sign export fullpubkey` - Export full public key addresses
- `sign import privkey` - Import private keys
- `sign sign` - Sign unsigned transactions (second and subsequent signatures)

## Transaction Execution Flow

### For Non-Multisig Addresses

```text
1. Watch Wallet: Create unsigned transaction
   └─> watch create deposit/payment/transfer
   └─> Transaction file (unsigned) is generated

2. Keygen Wallet: Execute first signature
   └─> keygen sign -file <tx_file>
   └─> Signed transaction file is generated

3. Watch Wallet: Send signed transaction
   └─> watch send -file <signed_tx_file>
   └─> Transaction is sent to blockchain and transaction ID is returned
```

### For Multisig Addresses

```text
1. Watch Wallet: Create unsigned transaction
   └─> watch create deposit/payment/transfer
   └─> Transaction file (unsigned) is generated

2. Keygen Wallet: Execute first signature
   └─> keygen sign -file <tx_file>
   └─> Transaction file with first signature is generated

3. Sign Wallet #1: Execute second signature
   └─> sign sign -file <tx_file_signed1>
   └─> Transaction file with second signature is generated

4. Sign Wallet #2: Execute third signature (if needed)
   └─> sign sign -file <tx_file_signed2>
   └─> Transaction file with third signature is generated
   └─> (Repeat as needed based on multisig configuration until required number of signatures)

5. Watch Wallet: Send signed transaction
   └─> watch send -file <fully_signed_tx_file>
   └─> Transaction is sent to blockchain and transaction ID is returned
```

## Key Generation Flow

Key generation flow for creating multisig addresses:

```text
1. Keygen Wallet: Generate seed for account
   └─> keygen create seed

2. Keygen Wallet: Generate HD key
   └─> keygen create hdkey --account <account>

3. Sign Wallet #1: Generate seed for authentication account
   └─> sign create seed

4. Sign Wallet #1: Generate HD key for authentication account
   └─> sign create hdkey

5. Sign Wallet #1: Export full public key
   └─> sign export fullpubkey
   └─> CSV file is generated

6. Keygen Wallet: Import full public key exported from Sign Wallet
   └─> keygen import full-pubkey -file <fullpubkey.csv>

7. Keygen Wallet: Generate multisig address
   └─> keygen create multisig --account <account>
   └─> Multisig address is generated

8. Keygen Wallet: Export public key address
   └─> keygen export address
   └─> CSV file is generated

9. Watch Wallet: Import address exported from Keygen Wallet
   └─> watch import address -file <address.csv>
   └─> Watch Wallet can now monitor the address
```

## Transaction Types

### Deposit

A transaction that consolidates coins sent to client addresses into a secure offline-managed address (cold wallet).

```bash
# Create unsigned transaction in Watch Wallet
watch create deposit

# Sign in Keygen Wallet
keygen sign -file <tx_file>

# Send in Watch Wallet
watch send -file <signed_tx_file>
```

### Payment

A transaction that sends coins to a specified address based on user withdrawal requests.

```bash
# Create unsigned transaction in Watch Wallet
watch create payment

# First signature in Keygen Wallet
keygen sign -file <tx_file>

# Second and subsequent signatures in Sign Wallet (for multisig)
sign sign -file <tx_file_signed1>

# Send in Watch Wallet
watch send -file <fully_signed_tx_file>
```

### Transfer

A transaction that transfers coins between internal accounts.

```bash
# Create unsigned transaction in Watch Wallet
watch create transfer --account1 <sender> --account2 <receiver> --amount <amount>

# Sign in Keygen Wallet
keygen sign -file <tx_file>

# Send in Watch Wallet
watch send -file <signed_tx_file>
```

## Security Design

- **Private Key Separation**: Watch Wallet does not hold any private keys and only manages public keys
- **Offline Operation**: Keygen Wallet and Sign Wallet operate in offline environments, isolated from networks
- **Multisig**: Multisig addresses requiring multiple signatures eliminate single points of failure
- **Role Separation**: Key generation, signing, and sending functions are separated into different wallets

## Supported Coins

- **Bitcoin (BTC)**
- **Bitcoin Cash (BCH)**
- **Ethereum (ETH)**
- **ERC-20 Token**
- **Ripple (XRP)**

All three wallet types are implemented for each coin.
