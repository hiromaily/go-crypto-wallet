## Known Issues and Workarounds

### BCH Node RPC Bug: Multisig `complete` Flag

#### Issue Description

Bitcoin Cash Node has a bug in the `signrawtransactionwithkey` RPC method where it incorrectly reports `complete=true` after adding a partial signature to multisig transactions:

| Multisig Type | Expected Behavior | Actual BCH Behavior | Impact |
|---------------|-------------------|---------------------|--------|
| 2-of-3 | `complete=false` after 1st sig, `complete=true` after 2nd sig | Same (correct) | None |
| 3-of-3 | `complete=false` after 1st and 2nd sig, `complete=true` after 3rd sig | `complete=true` after 2nd sig (incorrect) | Subsequent signers cannot add their signatures |

#### Root Cause

The BCH Node's `signrawtransactionwithkey` RPC method incorrectly evaluates transaction completeness for multisig transactions. When adding the 2nd signature to a 3-of-3 multisig transaction, it reports `complete=true` even though a 3rd signature is still required.

#### Symptoms

1. **Transaction File Type**: After 2nd signature in 3-of-3 multisig, transaction file is marked as `signed` instead of remaining `unsigned`
2. **Missing prevTx Metadata**: When marked as `signed`, prevTx metadata (TXID, vout, scriptPubKey, redeemScript, amount) is omitted from the transaction file
3. **Sign2 Wallet Failure**: Without prevTx metadata, Sign2 wallet cannot add the 3rd signature, resulting in error:

   ```
   Error: fail to sign transaction: fail to sign raw transaction with auth key:
   result of sign raw transaction includes error: Input not found or already spent
   ```

#### Workaround Implementation

**Issue #485** implemented the following workarounds:

##### 1. Accept Both `unsigned` and `signed` Transaction Types (Sign Wallet)

Modified `internal/application/usecase/sign/bch/sign_transaction.go` to accept both transaction types:

```go
// For BCH multisig, accept both "unsigned" and "signed" files since BCH may mark
// a transaction as "signed" even when it needs additional signatures for multisig
if fileType.TxType != domainTx.TxTypeUnsigned && fileType.TxType != domainTx.TxTypeSigned {
    return signusecase.SignTransactionOutput{}, fmt.Errorf(
        "invalid txType: %s (expected unsigned or signed)", fileType.TxType)
}
```

##### 2. Always Include prevTx Metadata for Multisig (Keygen and Sign Wallets)

Modified both keygen and sign wallets to always include prevTx metadata for multisig transactions, regardless of the `complete` flag:

**Keygen Wallet** (`internal/application/usecase/keygen/bch/sign_transaction.go`):

```go
var generatedFileName string
// For multisig transactions, always include prevTx metadata for subsequent signers,
// even if BCH incorrectly reports isSigned=true for partial signatures
isMultisig := u.multisigAccount != nil
if isSigned && !isMultisig {
    // For fully signed single-sig transactions, just write the hex
    generatedFileName, err = u.txFileRepo.WriteHexFile(basePath, signedHex)
} else {
    // For partially signed or multisig, include prevTx metadata for next signer
    content := u.formatSignedTxContent(signedHex, prevTxs)
    generatedFileName, err = u.txFileRepo.WriteHexFile(basePath, content)
}
```

**Sign Wallet** (`internal/application/usecase/sign/bch/sign_transaction.go`):

```go
var generatedFileName string
// For multisig transactions, always include prevTx metadata for subsequent signers,
// even if BCH incorrectly reports isSigned=true for partial signatures
isMultisig := u.multisigAccount != nil
if isSigned && !isMultisig {
    // For fully signed single-sig transactions, just write the hex
    generatedFileName, err = u.txFileRepo.WriteFile(path, signedHex)
} else {
    // For partially signed or multisig, include prevTx metadata for next signer
    content := u.formatSignedTxContent(signedHex, prevTxs)
    generatedFileName, err = u.txFileRepo.WriteFile(path, content)
}
```

##### 3. E2E Script Fixes

Fixed duplicate `--wallet` flags in E2E scripts:

- **Issue**: RPC host already included wallet name (e.g., `127.0.0.1:30332/wallet/sign1-p2`)
- **Additional `--wallet` flag** created invalid path (e.g., `127.0.0.1:30332/wallet/sign1-p2/wallet/sign1`)
- **Solution**: Removed all redundant `--wallet` flags from `e2e-p2-p2sh-2of3.sh` and `e2e-p3-p2sh-3of3.sh`

#### Testing Results

After implementing the workarounds:

| Pattern | Description | Status | Transaction ID |
|---------|-------------|--------|----------------|
| Pattern 2 | 2-of-3 Multisig | ✅ Pass | `62c4331beddc2279c1567e41a23cca16d0ed4e25b376578ee4a3980c3f6cb240` |
| Pattern 3 | 3-of-3 Multisig | ✅ Pass | `4dee6d3978176f02ce01a9b4dacbfc01e84f662d1a4ab8ab8e7bb00c3a413a0c` |

#### Impact on Development

- **Multisig Detection**: Use non-nil `multisigAccount` field to determine if transaction requires multisig handling
- **prevTx Metadata**: Always include for multisig, regardless of `complete` flag
- **Transaction Type Validation**: Sign wallets must accept both `unsigned` and `signed` transaction files
- **E2E Scripts**: Avoid redundant `--wallet` flags when RPC host already includes wallet name

#### Future Considerations

This workaround should be maintained until Bitcoin Cash Node fixes the `complete` flag behavior. Monitor BCH Node releases for:

- Fixes to `signrawtransactionwithkey` multisig evaluation
- Changes to transaction completeness logic
- Updates to RPC response format

#### Related Issues

- GitHub Issue #485: "Bitcoin Core wallet not loaded when importing private keys"
- Related commits:
  - `69a2df37` - Initial Pattern 2 fix
  - `1902e9b2` - Extended fix to Pattern 3

---
