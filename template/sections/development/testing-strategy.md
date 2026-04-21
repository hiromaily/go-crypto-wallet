## Testing Strategy

### Unit Tests

Location: `internal/application/usecase/*/btc/*_test.go`

**Current Approach:**

- Constructor tests verify use case instantiation
- Interface compliance tests verify correct interface implementation

**Example:**

```go
func TestNewSignTransactionUseCase(t *testing.T) {
    t.Run("creates use case successfully with nil dependencies", func(t *testing.T) {
        useCase := btc.NewSignTransactionUseCase(
            nil, // btc
            nil, // accountKeyRepo
            nil, // txFileRepo
            nil, // multisigAccount
            domainWallet.WalletTypeKeygen,
            "auth1",
        )
        assert.NotNil(t, useCase)
    })

    t.Run("returns correct interface type", func(t *testing.T) {
        useCase := btc.NewSignTransactionUseCase(...)
        assert.Implements(t, (*keygusecase.SignTransactionUseCase)(nil), useCase)
    })
}
```

### Integration Tests

**Requirements for Full Integration Tests:**

1. **Mock Bitcoin Client**
   - CreatePSBT
   - SignPSBTWithKey
   - FinalizePSBT
   - ExtractTransaction
   - IsPSBTComplete

2. **Mock Repositories**
   - TransactionFileRepository (read/write PSBT)
   - AccountKeyRepository (get keys)
   - AuthKeyRepository (get auth keys)
   - BTCTxRepository (database operations)

3. **Test Fixtures**
   - Sample PSBTs (unsigned, partially signed, fully signed)
   - Sample private keys (WIF format)
   - Sample transaction data

**Example Integration Test:**

```go
func TestSignTransactionUseCase_Integration(t *testing.T) {
    // Setup mocks
    mockBTC := &mockBitcoinClient{}
    mockKeyRepo := &mockAccountKeyRepository{}
    mockFileRepo := &mockTransactionFileRepository{}

    // Create use case
    useCase := btc.NewSignTransactionUseCase(
        mockBTC,
        mockKeyRepo,
        mockFileRepo,
        nil,
        domainWallet.WalletTypeKeygen,
        "auth1",
    )

    // Setup test data
    unsignedPSBT := loadTestPSBT("testdata/unsigned.psbt")
    mockFileRepo.On("ReadPSBTFile", mock.Anything).Return(unsignedPSBT, nil)
    mockKeyRepo.On("GetAll", mock.Anything, mock.Anything).Return(testKeys, nil)
    mockBTC.On("SignPSBTWithKey", mock.Anything, mock.Anything).Return(signedPSBT, true, nil)

    // Execute
    output, err := useCase.Sign(context.Background(), input)

    // Assert
    assert.NoError(t, err)
    assert.True(t, output.IsComplete)
    assert.NotEmpty(t, output.SignedData)
}
```

### End-to-End Tests

**Manual E2E Test on Testnet:**

```bash
# 1. Create unsigned PSBT
./watch create deposit --fee 0.00001

# 2. Sign with Keygen
./keygen sign --file deposit_*_unsigned_0_*.psbt

# 3. Broadcast
./watch send --file deposit_*_signed_1_*.psbt

# 4. Verify on blockchain
bitcoin-cli -testnet getrawtransaction <txid> 1
```

**Automated E2E Tests:**

See `docs/TESTING_STRATEGY.md` for comprehensive testing approach.

---
