## Database Schema

### Current Tables (Used by MuSig2)

#### account_key Table

Stores account keys including MuSig2 Taproot addresses:

```sql
CREATE TABLE account_key (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    account VARCHAR(255) NOT NULL,
    key_index INT NOT NULL,
    full_public_key TEXT NOT NULL,
    wallet_import_format TEXT,
    p2pkh_address VARCHAR(255),
    p2sh_segwit_address VARCHAR(255),
    bech32_address VARCHAR(255),
    multisig_address VARCHAR(255),     -- MuSig2 P2TR address stored here
    redeem_script TEXT,
    addr_status TINYINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_coin_account (coin, account),
    KEY idx_addr_status (addr_status)
);
```

**MuSig2 Usage**:

- `multisig_address`: Stores P2TR (Taproot) address (`bc1p...` or `tb1p...`)
- `addr_status`: Tracks address creation status
- `full_public_key`: Account's public key for aggregation

#### auth_fullpubkey Table

Stores public keys from Sign wallets (auth accounts):

```sql
CREATE TABLE auth_fullpubkey (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    auth_account VARCHAR(255) NOT NULL,
    full_public_key TEXT NOT NULL,
    p2pkh_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_coin_auth (coin, auth_account)
);
```

**MuSig2 Usage**:

- `full_public_key`: Sign wallet's public key for key aggregation
- `auth_account`: Identifies which Sign wallet (auth1, auth2, etc.)

#### auth_account_key Table

Stores Sign wallet's private keys for signing:

```sql
CREATE TABLE auth_account_key (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    auth_account VARCHAR(255) NOT NULL,
    full_public_key TEXT NOT NULL,
    wallet_import_format TEXT NOT NULL,
    p2pkh_address VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_coin_auth (coin, auth_account)
);
```

**MuSig2 Usage**:

- `wallet_import_format`: Private key for creating partial signatures
- Used by Sign wallets during Round 2 (signing)

### Future Tables (For Enhanced Nonce Management)

#### musig2_nonces Table (Proposed)

Track nonces to prevent reuse:

```sql
CREATE TABLE musig2_nonces (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    transaction_id BIGINT NOT NULL,
    signer_id VARCHAR(255) NOT NULL,
    nonce BINARY(66) NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_at TIMESTAMP NULL,
    UNIQUE KEY unique_nonce (nonce),
    UNIQUE KEY unique_tx_signer (transaction_id, signer_id),
    KEY idx_coin_tx (coin, transaction_id),
    KEY idx_used (used)
);
```

**Purpose**:

- Enforce nonce uniqueness at database level
- Track which nonces have been used
- Prevent nonce reuse attacks
- Audit trail for debugging

#### musig2_sessions Table (Proposed)

Track signing sessions:

```sql
CREATE TABLE musig2_sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL UNIQUE,
    coin VARCHAR(10) NOT NULL,
    transaction_id BIGINT NOT NULL,
    participant_count INT NOT NULL,
    nonces_collected INT DEFAULT 0,
    partial_sigs_collected INT DEFAULT 0,
    status ENUM('nonce_generation', 'signing', 'aggregating', 'completed', 'failed') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    KEY idx_coin_tx (coin, transaction_id),
    KEY idx_status (status)
);
```

**Purpose**:

- Track signing session progress
- Validate all nonces/signatures collected
- Monitor signing workflow state
- Debugging and audit trail

---
