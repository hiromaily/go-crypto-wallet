package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
)

type sendMPCTxDeps struct {
	mpcSigner   *ethapiamocks.MockMPCTransactionSigner
	txSender    *ethapiamocks.MockTxSender
	mpcFileRepo *storagemocks.MockMPCFileRepositorier
}

func newSendMPCTxDeps(t *testing.T) *sendMPCTxDeps {
	t.Helper()
	return &sendMPCTxDeps{
		mpcSigner:   ethapiamocks.NewMockMPCTransactionSigner(t),
		txSender:    ethapiamocks.NewMockTxSender(t),
		mpcFileRepo: storagemocks.NewMockMPCFileRepositorier(t),
	}
}

func newSendMPCTxUseCase(deps *sendMPCTxDeps) watchusecase.SendMPCTransactionUseCase {
	return watchusecaseeth.NewSendMPCTransactionUseCase(
		deps.mpcSigner, deps.txSender, deps.mpcFileRepo)
}

func TestSendMPCTransaction_ErrNotUnsigned(t *testing.T) {
	t.Parallel()

	deps := newSendMPCTxDeps(t)

	signedFile := &dtoeth.ETHMPCTransactionFile{
		Version:              1,
		TxType:               "signed",
		UUID:                 mpcUUIDStr,
		ActionType:           "payment",
		From:                 mpcFromAddr,
		To:                   mpcToAddr,
		Value:                "100000000000000000",
		GasLimit:             21000,
		ChainID:              31337,
		MaxFeePerGas:         "2000000000",
		MaxPriorityFeePerGas: "1000000000",
		TxHash:               mpcTestTxHash,
		RawTxHex:             mpcTestTxHex,
		SignedTxHex:          "0x1234signed",
		Threshold:            2,
		PartyIDs:             []string{"node-1", "node-2"},
	}
	deps.mpcFileRepo.EXPECT().
		ReadETHMPCJSONFile("/tmp/payment_mpc_test.json").
		Return(signedFile, nil).Once()

	uc := newSendMPCTxUseCase(deps)
	_, err := uc.Execute(context.Background(), watchusecase.SendMPCTransactionInput{
		FilePath:  "/tmp/payment_mpc_test.json",
		PeerAddrs: []string{"localhost:9001", "localhost:9002"},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, dtoeth.ErrNotUnsigned)
}

func TestSendMPCTransaction_ReadFileError(t *testing.T) {
	t.Parallel()

	deps := newSendMPCTxDeps(t)
	readErr := errors.New("file not found")

	deps.mpcFileRepo.EXPECT().
		ReadETHMPCJSONFile(mock.Anything).
		Return(nil, readErr).Once()

	uc := newSendMPCTxUseCase(deps)
	_, err := uc.Execute(context.Background(), watchusecase.SendMPCTransactionInput{
		FilePath:  "/tmp/missing.json",
		PeerAddrs: []string{"localhost:9001"},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, readErr)
}

func TestSendMPCTransaction_SenderMismatch(t *testing.T) {
	t.Parallel()

	const chainID = uint64(31337)
	// Build the tx signed with a test private key.
	privKey, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   new(big.Int).SetUint64(chainID),
		Nonce:     1,
		GasTipCap: big.NewInt(1e9),
		GasFeeCap: big.NewInt(2e9),
		Gas:       21000,
		Value:     big.NewInt(1e17),
	})

	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)

	// Sign the tx hash directly to get a valid 65-byte signature.
	signer := types.LatestSignerForChainID(new(big.Int).SetUint64(chainID))
	hash := signer.Hash(tx)
	sig65, err := ethcrypto.Sign(hash.Bytes(), privKey)
	require.NoError(t, err)

	// Adjust v from {0,1} to {27,28} to simulate what coordinator returns.
	sig65[64] += 27

	file := &dtoeth.ETHMPCTransactionFile{
		Version:              1,
		TxType:               "unsigned",
		UUID:                 mpcUUIDStr,
		ActionType:           "payment",
		From:                 mpcFromAddr, // different from the private key's address
		To:                   mpcToAddr,
		Value:                "100000000000000000",
		GasLimit:             21000,
		ChainID:              chainID,
		MaxFeePerGas:         "2000000000",
		MaxPriorityFeePerGas: "1000000000",
		TxHash:               tx.Hash().Hex(),
		RawTxHex:             *txHex,
		Threshold:            2,
		PartyIDs:             []string{"node-1", "node-2"},
	}

	deps := newSendMPCTxDeps(t)
	deps.mpcFileRepo.EXPECT().
		ReadETHMPCJSONFile(mock.Anything).
		Return(file, nil).Once()
	deps.mpcSigner.EXPECT().
		SignTransaction(mock.Anything, mock.Anything).
		Return(&apieth.MPCSigningResult{
			SessionID: mpcUUIDStr,
			Signature: sig65,
		}, nil).Once()

	uc := newSendMPCTxUseCase(deps)
	_, err = uc.Execute(context.Background(), watchusecase.SendMPCTransactionInput{
		FilePath:  "/tmp/test.json",
		PeerAddrs: []string{"localhost:9001", "localhost:9002"},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, watchusecaseeth.ErrSenderMismatch)
}

func TestSendMPCTransaction_SigningError(t *testing.T) {
	t.Parallel()

	privKey, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	_ = privKey

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID: big.NewInt(31337),
		Gas:     21000,
		Value:   big.NewInt(1e17),
	})
	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)

	file := &dtoeth.ETHMPCTransactionFile{
		Version:              1,
		TxType:               "unsigned",
		UUID:                 mpcUUIDStr,
		ActionType:           "payment",
		From:                 mpcFromAddr,
		To:                   mpcToAddr,
		Value:                "100000000000000000",
		GasLimit:             21000,
		ChainID:              31337,
		MaxFeePerGas:         "2000000000",
		MaxPriorityFeePerGas: "1000000000",
		TxHash:               tx.Hash().Hex(),
		RawTxHex:             *txHex,
		Threshold:            2,
		PartyIDs:             []string{"node-1", "node-2"},
	}

	sigErr := errors.New("TSS round timeout: session expired")
	deps := newSendMPCTxDeps(t)
	deps.mpcFileRepo.EXPECT().
		ReadETHMPCJSONFile(mock.Anything).
		Return(file, nil).Once()
	deps.mpcSigner.EXPECT().
		SignTransaction(mock.Anything, mock.Anything).
		Return(nil, sigErr).Once()

	uc := newSendMPCTxUseCase(deps)
	_, err = uc.Execute(context.Background(), watchusecase.SendMPCTransactionInput{
		FilePath:  "/tmp/test.json",
		PeerAddrs: []string{"localhost:9001"},
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, sigErr)
}
