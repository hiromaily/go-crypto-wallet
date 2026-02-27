package eth_test

// ETH file exchange integration tests: Watch → Keygen → Watch flow.
//
// These tests verify that the three use cases (Create, Keygen Sign, Send) correctly
// hand off data through JSON transaction files. All dependencies are mocked — no
// Ethereum node or database is required. The test exercises:
//
//   - EIP-1559 path: SupportsEIP1559 returns true → DynamicFeeTx file produced
//   - Legacy fallback: SupportsEIP1559 returns false → LegacyTx file produced
//   - File exchange: Watch create → Keygen sign → Watch send → Monitor TxTypeDone
//
// For tests that require a live Anvil node, see the sibling `//go:build integration`
// file in internal/integration_test/.

import (
	"context"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	persistencemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/persistence/mocks"
	coldmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

// ---------------------------------------------------------------------------
// Helpers shared across file-exchange tests
// ---------------------------------------------------------------------------

const (
	flowFromAddr     = "0xAbcdef1234567890AbCdEf1234567890AbCdEf12"
	flowToAddr       = "0xDeadBeef1234567890DeadBeef1234567890abcd"
	flowRawTxHex     = "0x02c0"
	flowSignedTxHex  = "0x02f8"
	flowTxID         = int64(7)
	flowUnsignedPath = "deposit_7_unsigned_0_1534744535097796209.json"
	flowSignedPath   = "deposit_7_signed_1_1534744535097796209.json"
	flowTxHash       = "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"
	flowUUID         = "deadbeef-dead-beef-dead-beefdeadbeef"
	flowChainID      = uint64(31337)
)

// newFlowAccountXpriv generates a valid BIP-44 ETH account-level xpriv for testing.
func newFlowAccountXpriv(t *testing.T) string {
	t.Helper()
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i + 1)
	}
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)
	purpose, err := master.Derive(hdkeychain.HardenedKeyStart + 44)
	require.NoError(t, err)
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 60)
	require.NoError(t, err)
	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0)
	require.NoError(t, err)
	return account.String()
}

// buildUnsignedEIP1559File returns a minimal unsigned EIP-1559 transaction file
// for use as the output of the Watch create step.
func buildUnsignedEIP1559File() *dtoeth.ETHTransactionFile {
	return &dtoeth.ETHTransactionFile{
		Version:              1,
		TxType:               string(domainTx.TxTypeUnsigned),
		EthTxType:            dtoeth.EthTxTypeEIP1559,
		UUID:                 flowUUID,
		ChainID:              flowChainID,
		Nonce:                0,
		From:                 flowFromAddr,
		To:                   flowToAddr,
		Value:                "999000000000000000",
		Gas:                  21000,
		MaxFeePerGas:         "20000000000",
		MaxPriorityFeePerGas: "2000000000",
		RawTxHex:             flowRawTxHex,
	}
}

// ---------------------------------------------------------------------------
// Test: EIP-1559 path — Watch creates DynamicFeeTx file, Keygen signs,
//       Watch sends and DB is updated.
// ---------------------------------------------------------------------------

// TestETHFileExchange_EIP1559Path tests the Watch→Keygen→Watch flow using
// EIP-1559 transaction type. The Watch wallet selects EIP-1559 when
// SupportsEIP1559 returns true, writes an unsigned file, and the Keygen
// wallet signs it to produce a signed file that Watch then broadcasts.
func TestETHFileExchange_EIP1559Path(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// ---- Mocks for Watch Create step ----
	createEthClient := ethapiamocks.NewMockERC20er(t)
	createUoW := persistencemocks.NewMockUnitOfWork(t)
	createTx := persistencemocks.NewMockTransaction(t)
	createAddrRepo := repomocks.NewMockAddressRepositorier(t)
	createTxRepo := repomocks.NewMockTxRepositorier(t)
	createTxDetailRepo := repomocks.NewMockETHDetailTXRepositorier(t)
	createPayReqRepo := repomocks.NewMockPaymentRequestRepositorier(t)
	createTxFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)

	// ---- Mocks for Watch Send step (separate client mock) ----
	sendEthClient := ethapiamocks.NewMockTxSender(t)
	sendTxDetailRepo := repomocks.NewMockETHDetailTXRepositorier(t)
	sendTxFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)

	// ---- Mocks for Keygen Sign step ----
	signTxSigner := ethapiamocks.NewMockTxSigner(t)
	signTxFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)
	signAccountKeyRepo := coldmocks.NewMockETHAccountKeyRepositorier(t)

	// ============================================================
	// Step 1: Watch Create → writes unsigned EIP-1559 file
	// ============================================================
	// addrRepo.GetAll returns one client address
	createAddrRepo.EXPECT().
		GetAll(domainAccount.AccountTypeClient).
		Return([]*domainAddress.Address{{WalletAddress: flowFromAddr}}, nil)

	// GetBalance returns 1 ETH for the client address
	createEthClient.EXPECT().
		GetBalance(ctx, flowFromAddr, domainETH.QuantityTagLatest).
		Return(new(big.Int).SetUint64(1_000_000_000_000_000_000), nil)

	// GetOneUnAllocated returns deposit address
	createAddrRepo.EXPECT().
		GetOneUnAllocated(domainAccount.AccountTypeDeposit).
		Return(&domainAddress.Address{WalletAddress: flowToAddr}, nil)

	// EIP-1559 path
	createEthClient.EXPECT().
		SupportsEIP1559(ctx).
		Return(true)
	createEthClient.EXPECT().
		CreateRawTransactionEIP1559(ctx, flowFromAddr, flowToAddr, uint64(0), 0).
		Return(
			&domainETH.RawTx{
				UUID:  flowUUID,
				From:  flowFromAddr,
				To:    flowToAddr,
				Value: *new(big.Int).SetUint64(999_000_000_000_000_000),
				Nonce: 0,
				TxHex: flowRawTxHex,
			},
			&apieth.TxCreateParams{
				UUID:                 flowUUID,
				FromAddress:          flowFromAddr,
				ToAddress:            flowToAddr,
				Amount:               999_000_000_000_000_000,
				Fee:                  1_000_000_000_000_000,
				GasLimit:             21000,
				Nonce:                0,
				EthTxType:            dtoeth.EthTxTypeEIP1559,
				ChainID:              flowChainID,
				MaxFeePerGas:         20_000_000_000,
				MaxPriorityFeePerGas: 2_000_000_000,
			},
			nil,
		)

	// DB transaction: Begin → WithTransaction → Insert → Commit
	createUoW.EXPECT().
		Begin(ctx).
		Return(createTx, nil)
	createTxRepo.EXPECT().
		WithTransaction(createTx).
		Return(createTxRepo, nil)
	createTxDetailRepo.EXPECT().
		WithTransaction(createTx).
		Return(createTxDetailRepo, nil)
	createPayReqRepo.EXPECT().
		WithTransaction(createTx).
		Return(createPayReqRepo, nil)
	createTxRepo.EXPECT().
		InsertUnsignedTx(domainTx.ActionTypeDeposit).
		Return(flowTxID, nil)
	createTxDetailRepo.EXPECT().
		InsertBulk(mock.AnythingOfType("[]*eth.ETHDetailTx")).
		Return(nil)
	createTx.EXPECT().
		Commit().
		Return(nil)

	// File repo
	createTxFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, flowTxID, 0).
		Return(flowUnsignedPath)
	createTxFileRepo.EXPECT().
		WriteETHJSONFile(
			flowUnsignedPath,
			mock.MatchedBy(func(f *dtoeth.ETHTransactionFile) bool {
				return f.EthTxType == dtoeth.EthTxTypeEIP1559 &&
					f.TxType == string(domainTx.TxTypeUnsigned) &&
					f.MaxFeePerGas == "20000000000" &&
					f.GasPrice == ""
			}),
		).
		Return(flowUnsignedPath, nil)

	createUC := watchusecaseeth.NewCreateTransactionUseCase(
		createEthClient,
		createUoW,
		createAddrRepo,
		createTxRepo,
		createTxDetailRepo,
		createPayReqRepo,
		createTxFileRepo,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
	)

	createOut, err := createUC.Execute(ctx, watchusecase.CreateTransactionInput{
		ActionType: string(domainTx.ActionTypeDeposit),
	})
	require.NoError(t, err)
	require.Equal(t, flowUnsignedPath, createOut.FileName)

	// ============================================================
	// Step 2: Keygen Sign → signs unsigned file → writes signed file
	// ============================================================
	xpriv := newFlowAccountXpriv(t)

	ethKey, err := domainETH.NewETHAccountKey(
		domainAccount.AccountTypeDeposit,
		flowFromAddr,
		"0x04aabbcc",
		"someprivkey",
		0,
	)
	require.NoError(t, err)
	ethKey.SetAccountExtendedPrivkey(xpriv)

	signTxFileRepo.EXPECT().
		ValidateFilePath(flowUnsignedPath, domainTx.TxTypeUnsigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, flowTxID, 0, nil)
	signTxFileRepo.EXPECT().
		ReadETHJSONFile(flowUnsignedPath).
		Return(buildUnsignedEIP1559File(), nil)
	signAccountKeyRepo.EXPECT().
		GetByAddress(flowFromAddr).
		Return(ethKey, nil)
	signTxSigner.EXPECT().
		SignTxWithPrivateKey(
			mock.MatchedBy(func(rawTx *domainETH.RawTx) bool {
				return rawTx.TxHex == flowRawTxHex
			}),
			mock.Anything,
			mock.MatchedBy(func(chainID *big.Int) bool {
				return chainID.Uint64() == flowChainID
			}),
		).
		Return(&domainETH.RawTx{
			From:  flowFromAddr,
			To:    flowToAddr,
			TxHex: flowSignedTxHex,
		}, nil)
	signTxFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, flowTxID, 1).
		Return(flowSignedPath)
	signTxFileRepo.EXPECT().
		WriteETHJSONFile(
			flowSignedPath,
			mock.MatchedBy(func(f *dtoeth.ETHTransactionFile) bool {
				return f.TxType == string(domainTx.TxTypeSigned) &&
					f.SignedTxHex == flowSignedTxHex
			}),
		).
		Return(flowSignedPath, nil)

	signUC := keygenusecaseeth.NewSignTransactionUseCase(
		signTxSigner,
		signTxFileRepo,
		signAccountKeyRepo,
	)

	signOut, err := signUC.Sign(ctx, keygenusecase.SignTransactionInput{
		FilePath: flowUnsignedPath,
	})
	require.NoError(t, err)
	require.Equal(t, flowSignedPath, signOut.FilePath)
	require.True(t, signOut.IsDone)

	// ============================================================
	// Step 3: Watch Send → reads signed file → broadcasts → DB updated
	// ============================================================
	sendTxFileRepo.EXPECT().
		ValidateFilePath(flowSignedPath, domainTx.TxTypeSigned).
		Return(domainTx.ActionTypeDeposit, domainTx.TxTypeSigned, flowTxID, 1, nil)

	signedFile := buildUnsignedEIP1559File()
	signedFile.TxType = string(domainTx.TxTypeSigned)
	signedFile.SignedTxHex = flowSignedTxHex

	sendTxFileRepo.EXPECT().
		ReadETHJSONFile(flowSignedPath).
		Return(signedFile, nil)
	sendEthClient.EXPECT().
		SendSignedRawTransaction(ctx, flowSignedTxHex).
		Return(flowTxHash, nil)
	sendTxDetailRepo.EXPECT().
		UpdateAfterTxSent(flowUUID, domainTx.TxTypeSent, flowSignedTxHex, flowTxHash).
		Return(int64(1), nil)

	sendUC := watchusecaseeth.NewSendTransactionUseCase(
		sendEthClient,
		sendTxDetailRepo,
		sendTxFileRepo,
	)

	sendOut, err := sendUC.Execute(ctx, watchusecase.SendTransactionInput{
		FilePath: flowSignedPath,
	})
	require.NoError(t, err)
	require.Equal(t, flowTxHash, sendOut.TxID)
}

// ---------------------------------------------------------------------------
// Test: Legacy fallback path — SupportsEIP1559 returns false
// ---------------------------------------------------------------------------

// TestETHFileExchange_LegacyFallback verifies that when the Ethereum node does not
// support EIP-1559, the Watch wallet falls back to creating a legacy (Type 0)
// transaction file with GasPrice instead of MaxFeePerGas / MaxPriorityFeePerGas.
func TestETHFileExchange_LegacyFallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	createEthClient := ethapiamocks.NewMockERC20er(t)
	createUoW := persistencemocks.NewMockUnitOfWork(t)
	createTx := persistencemocks.NewMockTransaction(t)
	createAddrRepo := repomocks.NewMockAddressRepositorier(t)
	createTxRepo := repomocks.NewMockTxRepositorier(t)
	createTxDetailRepo := repomocks.NewMockETHDetailTXRepositorier(t)
	createPayReqRepo := repomocks.NewMockPaymentRequestRepositorier(t)
	createTxFileRepo := storagemocks.NewMockTransactionFileRepositorier(t)

	createAddrRepo.EXPECT().
		GetAll(domainAccount.AccountTypeClient).
		Return([]*domainAddress.Address{{WalletAddress: flowFromAddr}}, nil)
	createEthClient.EXPECT().
		GetBalance(ctx, flowFromAddr, domainETH.QuantityTagLatest).
		Return(new(big.Int).SetUint64(1_000_000_000_000_000_000), nil)
	createAddrRepo.EXPECT().
		GetOneUnAllocated(domainAccount.AccountTypeDeposit).
		Return(&domainAddress.Address{WalletAddress: flowToAddr}, nil)

	// Legacy path: SupportsEIP1559 returns false
	createEthClient.EXPECT().
		SupportsEIP1559(ctx).
		Return(false)
	createEthClient.EXPECT().
		CreateRawTransaction(ctx, flowFromAddr, flowToAddr, uint64(0), 0).
		Return(
			&domainETH.RawTx{
				UUID:  flowUUID,
				From:  flowFromAddr,
				To:    flowToAddr,
				Value: *new(big.Int).SetUint64(999_000_000_000_000_000),
				Nonce: 0,
				TxHex: flowRawTxHex,
			},
			&apieth.TxCreateParams{
				UUID:        flowUUID,
				FromAddress: flowFromAddr,
				ToAddress:   flowToAddr,
				Amount:      999_000_000_000_000_000,
				Fee:         1_000_000_000_000_000,
				GasLimit:    21000,
				Nonce:       0,
				EthTxType:   dtoeth.EthTxTypeLegacy,
				ChainID:     flowChainID,
				GasPrice:    20_000_000_000,
			},
			nil,
		)

	createUoW.EXPECT().Begin(ctx).Return(createTx, nil)
	createTxRepo.EXPECT().WithTransaction(createTx).Return(createTxRepo, nil)
	createTxDetailRepo.EXPECT().WithTransaction(createTx).Return(createTxDetailRepo, nil)
	createPayReqRepo.EXPECT().WithTransaction(createTx).Return(createPayReqRepo, nil)
	createTxRepo.EXPECT().InsertUnsignedTx(domainTx.ActionTypeDeposit).Return(flowTxID, nil)
	createTxDetailRepo.EXPECT().
		InsertBulk(mock.AnythingOfType("[]*eth.ETHDetailTx")).
		Return(nil)
	createTx.EXPECT().Commit().Return(nil)

	createTxFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, flowTxID, 0).
		Return(flowUnsignedPath)
	createTxFileRepo.EXPECT().
		WriteETHJSONFile(
			flowUnsignedPath,
			mock.MatchedBy(func(f *dtoeth.ETHTransactionFile) bool {
				// Must use legacy fields, not EIP-1559 fields
				return f.EthTxType == dtoeth.EthTxTypeLegacy &&
					f.GasPrice == "20000000000" &&
					f.MaxFeePerGas == "" &&
					f.MaxPriorityFeePerGas == ""
			}),
		).
		Return(flowUnsignedPath, nil)

	createUC := watchusecaseeth.NewCreateTransactionUseCase(
		createEthClient,
		createUoW,
		createAddrRepo,
		createTxRepo,
		createTxDetailRepo,
		createPayReqRepo,
		createTxFileRepo,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
	)

	createOut, err := createUC.Execute(ctx, watchusecase.CreateTransactionInput{
		ActionType: string(domainTx.ActionTypeDeposit),
	})
	require.NoError(t, err)
	require.Equal(t, flowUnsignedPath, createOut.FileName)
}

// ---------------------------------------------------------------------------
// Test: TxTypeDone status transition via Monitor after Send
// ---------------------------------------------------------------------------

// TestETHFileExchange_TxTypeDoneTransition tests that after Watch sends a
// transaction, the Monitor use case correctly transitions the status to
// TxTypeDone once confirmations reach the threshold.
func TestETHFileExchange_TxTypeDoneTransition(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	monitorEthClient := ethapiamocks.NewMockEtherTxMonitor(t)
	monitorAddrRepo := repomocks.NewMockAddressRepositorier(t)
	monitorTxDetailRepo := repomocks.NewMockETHDetailTXRepositorier(t)

	// Sent transaction hash (as would be stored after Watch send step)
	monitorTxDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{flowTxHash}, nil)

	// Receipt shows success
	monitorEthClient.EXPECT().
		GetTxReceipt(ctx, flowTxHash).
		Return(&domainETH.TransactionReceipt{
			TransactionHash: flowTxHash,
			BlockNumber:     100,
			From:            flowFromAddr,
			To:              flowToAddr,
			Status:          1,
		}, nil)

	// Confirmations meet threshold (8 >= 6)
	monitorEthClient.EXPECT().
		GetConfirmation(ctx, flowTxHash).
		Return(uint64(8), nil)

	// Transition to TxTypeDone
	monitorTxDetailRepo.EXPECT().
		UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, flowTxHash).
		Return(int64(1), nil)

	// Free the sender address
	monitorAddrRepo.EXPECT().
		UpdateIsAllocated(false, flowFromAddr).
		Return(int64(1), nil)

	monitorUC := watchusecaseeth.NewMonitorTransactionUseCase(
		monitorEthClient,
		monitorAddrRepo,
		monitorTxDetailRepo,
		uint64(6), // confirmation threshold
	)

	err := monitorUC.UpdateTxStatus(ctx)
	require.NoError(t, err)
}
