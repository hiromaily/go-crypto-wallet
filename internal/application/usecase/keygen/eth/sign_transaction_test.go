package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

const (
	testFromAddr     = "0xAbcdef1234567890AbCdEf1234567890AbCdEf12"
	testToAddr       = "0xDeadBeef1234567890DeadBeef1234567890abcd"
	testRawTxHex     = "0x02c0" // minimal EIP-2718 encoded tx for unit testing (mocked signer)
	testSignedTxHex  = "0x02f8" // minimal signed tx hex for unit testing (mocked signer)
	testUnsignedPath = "deposit_8_unsigned_0_1534744535097796209.json"
	testSignedPath   = "deposit_8_signed_1_1534744535097796209.json"
	testChainID      = uint64(1337)
)

// newTestAccountXprivForSign generates a valid BIP-32 extended private key at account level for testing.
func newTestAccountXprivForSign(t *testing.T) string {
	t.Helper()
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i + 1)
	}
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	// Derive m/44'/60'/0' (BIP-44 ETH account 0)
	purpose, err := master.Derive(hdkeychain.HardenedKeyStart + 44)
	require.NoError(t, err)
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 60)
	require.NoError(t, err)
	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0)
	require.NoError(t, err)

	return account.String()
}

func newTestUnsignedTxFile() *dtoeth.ETHTransactionFile {
	return &dtoeth.ETHTransactionFile{
		Version:              1,
		TxType:               string(domainTx.TxTypeUnsigned),
		EthTxType:            2,
		ChainID:              testChainID,
		Nonce:                0,
		From:                 testFromAddr,
		To:                   testToAddr,
		Value:                "1000000000000000000",
		Gas:                  21000,
		MaxFeePerGas:         "20000000000",
		MaxPriorityFeePerGas: "1000000000",
		RawTxHex:             testRawTxHex,
	}
}

func newTestSignedRawTx() *domainETH.RawTx {
	return &domainETH.RawTx{
		From:  testFromAddr,
		To:    testToAddr,
		Value: *new(big.Int).SetUint64(1000000000000000000),
		TxHex: testSignedTxHex,
		Hash:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
	}
}

func TestSignTransactionUseCase_Sign_Success(t *testing.T) {
	t.Parallel()

	xpriv := newTestAccountXprivForSign(t)

	ethKey, err := domainETH.NewETHAccountKey(
		domainAccount.AccountTypeDeposit,
		testFromAddr,
		"0x04aabbcc",
		"someprivkey",
		0, // idx = 0, so derivation path: m/44'/60'/0'/0/0
	)
	require.NoError(t, err)
	ethKey.SetAccountExtendedPrivkey(xpriv)

	txSigner := ethapiamocks.NewMockTxSigner(t)
	txFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)
	accountKeyRepo := coldmocks.NewMockETHAccountKeyRepositorier(t)

	txFileRepo.EXPECT().
		ValidateFilePath(testUnsignedPath, domainTx.TxTypeUnsigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, int64(8), 0, nil)

	txFileRepo.EXPECT().
		ReadETHJSONFile(testUnsignedPath).
		Return(newTestUnsignedTxFile(), nil)

	accountKeyRepo.EXPECT().
		GetByAddress(testFromAddr).
		Return(ethKey, nil)

	txSigner.EXPECT().
		SignTxWithPrivateKey(
			mock.MatchedBy(func(rawTx *domainETH.RawTx) bool {
				return rawTx.TxHex == testRawTxHex
			}),
			mock.Anything, // derived *ecdsa.PrivateKey
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Uint64() == testChainID
			}),
		).
		Return(newTestSignedRawTx(), nil)

	txFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, int64(8), 1).
		Return(testSignedPath)

	txFileRepo.EXPECT().
		WriteETHJSONFile(
			testSignedPath,
			mock.MatchedBy(func(f *dtoeth.ETHTransactionFile) bool {
				return f.TxType == string(domainTx.TxTypeSigned) && f.SignedTxHex == testSignedTxHex
			}),
		).
		Return(testSignedPath, nil)

	useCase := keygenusecaseeth.NewSignTransactionUseCase(
		txSigner,
		txFileRepo,
		accountKeyRepo,
	)

	output, err := useCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: testUnsignedPath,
	})

	require.NoError(t, err)
	require.Equal(t, testSignedPath, output.FilePath)
	require.True(t, output.IsDone)
	require.Equal(t, 1, output.SignedCount)
}

func TestSignTransactionUseCase_Sign_InvalidFilePath(t *testing.T) {
	t.Parallel()

	txSigner := ethapiamocks.NewMockTxSigner(t)
	txFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)
	accountKeyRepo := coldmocks.NewMockETHAccountKeyRepositorier(t)

	txFileRepo.EXPECT().
		ValidateFilePath("invalid_path.json", domainTx.TxTypeUnsigned).
		Return(domainTx.ActionType(""), domainTx.TxType(""), int64(0), 0, errors.New("invalid file path"))

	useCase := keygenusecaseeth.NewSignTransactionUseCase(txSigner, txFileRepo, accountKeyRepo)

	_, err := useCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: "invalid_path.json",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid file path")
}

func TestSignTransactionUseCase_Sign_ReadFileError(t *testing.T) {
	t.Parallel()

	txSigner := ethapiamocks.NewMockTxSigner(t)
	txFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)
	accountKeyRepo := coldmocks.NewMockETHAccountKeyRepositorier(t)

	txFileRepo.EXPECT().
		ValidateFilePath(testUnsignedPath, domainTx.TxTypeUnsigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, int64(8), 0, nil)

	txFileRepo.EXPECT().
		ReadETHJSONFile(testUnsignedPath).
		Return(nil, errors.New("file not found"))

	useCase := keygenusecaseeth.NewSignTransactionUseCase(txSigner, txFileRepo, accountKeyRepo)

	_, err := useCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: testUnsignedPath,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to read ETH transaction file")
}

func TestSignTransactionUseCase_Sign_KeyNotFound(t *testing.T) {
	t.Parallel()

	txSigner := ethapiamocks.NewMockTxSigner(t)
	txFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)
	accountKeyRepo := coldmocks.NewMockETHAccountKeyRepositorier(t)

	txFileRepo.EXPECT().
		ValidateFilePath(testUnsignedPath, domainTx.TxTypeUnsigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, int64(8), 0, nil)

	txFileRepo.EXPECT().
		ReadETHJSONFile(testUnsignedPath).
		Return(newTestUnsignedTxFile(), nil)

	accountKeyRepo.EXPECT().
		GetByAddress(testFromAddr).
		Return(nil, errors.New("key not found"))

	useCase := keygenusecaseeth.NewSignTransactionUseCase(txSigner, txFileRepo, accountKeyRepo)

	_, err := useCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: testUnsignedPath,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "key not found for address")
}

func TestSignTransactionUseCase_Sign_MissingAccountXpriv(t *testing.T) {
	t.Parallel()

	// Account key without accountXpriv set
	ethKey, err := domainETH.NewETHAccountKey(
		domainAccount.AccountTypeDeposit,
		testFromAddr,
		"0x04aabbcc",
		"someprivkey",
		0,
	)
	require.NoError(t, err)
	// Do NOT call ethKey.SetAccountExtendedPrivkey(...)

	txSigner := ethapiamocks.NewMockTxSigner(t)
	txFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)
	accountKeyRepo := coldmocks.NewMockETHAccountKeyRepositorier(t)

	txFileRepo.EXPECT().
		ValidateFilePath(testUnsignedPath, domainTx.TxTypeUnsigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, int64(8), 0, nil)

	txFileRepo.EXPECT().
		ReadETHJSONFile(testUnsignedPath).
		Return(newTestUnsignedTxFile(), nil)

	accountKeyRepo.EXPECT().
		GetByAddress(testFromAddr).
		Return(ethKey, nil)

	useCase := keygenusecaseeth.NewSignTransactionUseCase(txSigner, txFileRepo, accountKeyRepo)

	_, err = useCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: testUnsignedPath,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "accountXpriv not found")
}
