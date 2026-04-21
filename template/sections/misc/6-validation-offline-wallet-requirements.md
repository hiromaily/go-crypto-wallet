## 6. Validation: Offline Wallet Requirements

### 6.1 Keygen Wallet (Offline)

**Requirements**:

- ✅ Read PSBT files from filesystem
- ✅ Parse PSBT without network access
- ✅ Sign PSBT using local private keys
- ✅ Write signed PSBT back to filesystem
- ✅ No Bitcoin Core RPC required

**btcd Package Support**:

- ✅ `psbt.NewFromRawBytes()` - Parse from base64
- ✅ `updater.Sign()` - Sign with private keys
- ✅ `packet.Serialize()` - Serialize to base64
- ✅ All operations local, no network calls

**Verdict**: ✅ **Fully Compatible**

### 6.2 Sign Wallet (Offline)

**Requirements**:

- ✅ Read partially signed PSBT files
- ✅ Parse PSBT with existing signatures
- ✅ Add second signature offline
- ✅ Write fully signed PSBT
- ✅ No Bitcoin Core RPC required

**btcd Package Support**:

- ✅ `psbt.NewFromRawBytes()` - Parse partially signed PSBT
- ✅ `updater.Sign()` - Add additional signatures
- ✅ `packet.IsComplete()` - Check completion
- ✅ All operations local

**Verdict**: ✅ **Fully Compatible**

### 6.3 Offline Operation Validation

| Operation | Keygen | Sign | Implementation |
|-----------|--------|------|----------------|
| Read PSBT file | ✅ | ✅ | `os.ReadFile()` |
| Parse base64 | ✅ | ✅ | `base64.StdEncoding.DecodeString()` |
| Parse PSBT | ✅ | ✅ | `psbt.NewFromRawBytes()` |
| Get private keys | ✅ | ✅ | Local database (SQLite) |
| Create signatures | ✅ | ✅ | btcd crypto functions |
| Add signatures | ✅ | ✅ | `updater.Sign()` |
| Serialize PSBT | ✅ | ✅ | `packet.Serialize()` |
| Encode base64 | ✅ | ✅ | `base64.StdEncoding.EncodeToString()` |
| Write PSBT file | ✅ | ✅ | `os.WriteFile()` |

**Verdict**: ✅ **All operations supported offline**

---
