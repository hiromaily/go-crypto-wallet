### Wallet Types Understanding

The project implements three types of wallets, each with specific roles:

#### Watch Wallet

- **Status**: Online
- **Keys**: Public keys only (no private keys)
- **Purpose**: Creates and sends transactions
- **Security**: Lower risk (no private keys exposed to online systems)
- **Use Cases**:
  - Monitor addresses and balances
  - Create unsigned transactions
  - Broadcast signed transactions
  - Query blockchain state

#### Keygen Wallet

- **Status**: Offline
- **Keys**: Generates and stores private keys
- **Purpose**: Generates keys, provides first signature for multisig
- **Security**: Highest security (offline, air-gapped if possible)
- **Use Cases**:
  - Generate new key pairs
  - Create multisig addresses
  - Provide first signature in multisig transactions
  - Export public keys to watch wallet

#### Sign Wallet

- **Status**: Offline
- **Keys**: Stores private keys for signing
- **Purpose**: Provides second and subsequent signatures for multisig
- **Security**: High security (offline)
- **Use Cases**:
  - Import unsigned transactions
  - Sign transactions with stored keys
  - Complete multisig signing process
  - Export signed transactions
