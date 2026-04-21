## Use Case Layer

### Watch Wallet Use Cases

#### CreateTransactionUseCase

Location: `internal/application/usecase/watch/btc/create_transaction.go`

**Responsibility:** Create unsigned PSBT for transactions.

**Key Method:**

```go
func (u *createTransactionUseCase) Execute(
    ctx context.Context,
    input watchusecase.CreateTransactionInput,
) (watchusecase.CreateTransactionOutput, error)
```

**PSBT Flow:**

1. Select UTXOs for inputs
2. Calculate outputs (recipient + change)
3. Create `wire.MsgTx`
4. Get previous transaction data for inputs
5. Call `btcClient.CreatePSBT(msgTx, prevTxs)`
6. Write PSBT to file
7. Store transaction metadata in database

**Code Example:**

```go
// Create transaction
msgTx, err := u.btcClient.CreateRawTransaction(inputs, outputs)

// Get previous transaction data
previousTxs, err := u.getPreviousTransactions(inputs)

// Create PSBT
psbtBase64, err := u.btcClient.CreatePSBT(msgTx, previousTxs.PrevTxs)

// Write PSBT file
path := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeUnsigned, txID, 0)
generatedFileName, err := u.txFileRepo.WritePSBTFile(path, psbtBase64)

return watchusecase.CreateTransactionOutput{
    TransactionHex: psbtBase64,
    FileName:       generatedFileName,
}, nil
```

#### SendTransactionUseCase

Location: `internal/application/usecase/watch/btc/send_transaction.go`

**Responsibility:** Finalize and broadcast fully signed PSBT.

**Key Method:**

```go
func (u *sendTransactionUseCase) Execute(
    ctx context.Context,
    input watchusecase.SendTransactionInput,
) (watchusecase.SendTransactionOutput, error)
```

**PSBT Flow:**

1. Detect file format (PSBT vs legacy)
2. For PSBT: Read PSBT file
3. Validate PSBT is fully signed
4. Finalize PSBT
5. Extract transaction
6. Convert to hex
7. Broadcast transaction
8. Update database

**Code Example:**

```go
func (u *sendTransactionUseCase) processPSBTFile(filePath string) (string, error) {
    // Read PSBT
    psbtBase64, err := u.txFileRepo.ReadPSBTFile(filePath)

    // Validate fully signed
    isComplete, err := u.btcClient.IsPSBTComplete(psbtBase64)
    if !isComplete {
        return "", errors.New("PSBT is not fully signed")
    }

    // Finalize PSBT
    finalizedPSBT, err := u.btcClient.FinalizePSBT(psbtBase64)

    // Extract transaction
    msgTx, err := u.btcClient.ExtractTransaction(finalizedPSBT)

    // Convert to hex
    hexTx, err := u.btcClient.ToHex(msgTx)

    return hexTx, nil
}
```

### Keygen Wallet Use Cases

#### SignTransactionUseCase (Keygen)

Location: `internal/application/usecase/keygen/btc/sign_transaction.go`

**Responsibility:** Add first signature to PSBT (offline).

**PSBT Flow:**

1. Read unsigned PSBT
2. Determine sender account
3. Get account private keys
4. Sign PSBT with keys (offline, no RPC)
5. Write partially/fully signed PSBT

**Code Example:**

```go
func (u *signTransactionUseCase) signMultisigPSBT(
    psbtBase64 string,
    senderAccount domainAccount.AccountType,
) (string, bool, error) {
    // Get account keys
    accountKeys, err := u.accountKeyRepo.GetAll(senderAccount, 0)

    // Extract WIFs from keys
    wifs := extractWIFs(accountKeys)

    // Sign PSBT offline (no Bitcoin Core RPC)
    signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(psbtBase64, wifs)

    return signedPSBT, isSigned, nil
}
```

### Sign Wallet Use Cases

#### SignTransactionUseCase (Sign)

Location: `internal/application/usecase/sign/btc/sign_transaction.go`

**Responsibility:** Add second+ signature to PSBT (offline).

**PSBT Flow:**

1. Read partially signed PSBT
2. Get auth private key
3. Sign PSBT with auth key (offline)
4. Write fully signed PSBT

**Code Example:**

```go
func (u *signTransactionUseCase) signMultisigPSBT(
    psbtBase64 string,
) (string, bool, error) {
    // Get auth key (explicit authType)
    authKey, err := u.authKeyRepo.GetOne(u.authType)

    // Sign PSBT offline
    signedPSBT, isSigned, err := u.btc.SignPSBTWithKey(
        psbtBase64,
        []string{authKey.WalletImportFormat},
    )

    return signedPSBT, isSigned, nil
}
```

---
