## Signing Mechanism

### Signer Types

| Signer | Transaction Type | Description |
|--------|-----------------|-------------|
| `types.HomesteadSigner` | Legacy | Original signer (no replay protection) |
| `types.NewEIP155Signer(chainID)` | Legacy | Replay protection (EIP-155) |
| `types.NewLondonSigner(chainID)` | Legacy + EIP-1559 | Used in this system |

### Signing Flow

```go
// 1. Decode hex transaction
tx := new(types.Transaction)
rlp.DecodeBytes(rawBytes, tx)

// 2. Get private key from keystore
privateKey := keystore.GetPrivKey(address, password)

// 3. Determine chain ID
// netID 1 → chainID 1 (mainnet)
// other → chainID 4 (testnet/regtest)

// 4. Create London signer (supports both legacy and EIP-1559)
signer := types.NewLondonSigner(chainID)

// 5. Sign transaction
signedTx, err := types.SignTx(tx, signer, privateKey)

// 6. Verify sender address
sender, err := types.Sender(signer, signedTx)

// 7. Encode to hex for file output
```

### Key Storage

Private keys are stored in the **go-ethereum keystore** format on the local filesystem:

```
File: UTC--{timestamp}--{address}
Path: Configurable (default: ./data/keystore/)
Encryption: scrypt (StandardScryptN, StandardScryptP)
Password: Required to decrypt (set in configuration)
```

**Key Import (Keygen Wallet):**

```go
// Uses local filesystem keystore (compatible with Anvil)
// Does NOT use personal_importRawKey RPC
ks := keystore.NewKeyStore(keyDir, keystore.StandardScryptN, keystore.StandardScryptP)
account, err := ks.ImportECDSA(ecdsaPrivKey, password)
```

---
