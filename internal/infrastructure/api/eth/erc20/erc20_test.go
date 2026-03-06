package erc20_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	apierc20impl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/erc20"
)

// stubEthNode is a minimal test double for the eth node provider (ERC20Operator).
// It records method calls for verification in TDD tests.
type stubEthNode struct {
	tipCap         *big.Int
	tipCapErr      error
	blockNumber    *big.Int
	blockNumberErr error
	blockInfo      *dtoeth.BlockInfo
	blockInfoErr   error
	txCount        *big.Int
	txCountErr     error

	// call trackers
	suggestGasTipCapCalled bool
	blockNumberCalled      bool
	getBlockByNumberCalled bool
}

func (s *stubEthNode) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	s.suggestGasTipCapCalled = true
	return s.tipCap, s.tipCapErr
}

func (s *stubEthNode) BlockNumber(_ context.Context) (*big.Int, error) {
	s.blockNumberCalled = true
	return s.blockNumber, s.blockNumberErr
}

func (s *stubEthNode) GetBlockByNumber(_ context.Context, _ uint64) (*dtoeth.BlockInfo, error) {
	s.getBlockByNumberCalled = true
	return s.blockInfo, s.blockInfoErr
}

func (s *stubEthNode) GetTransactionCount(_ context.Context, _ string, _ domainETH.QuantityTag) (*big.Int, error) {
	return s.txCount, s.txCountErr
}

// ---------------------------------------------------------------------------
// Task 3.1 – interface / constructor shape
// ---------------------------------------------------------------------------

// TestERC20_ImplementsERC20er verifies that NewERC20 returns a value satisfying ERC20er.
// The compile-time assertion (var _ apieth.ERC20er = (*erc20)(nil)) lives in erc20.go.
func TestERC20_ImplementsERC20er(t *testing.T) {
	t.Parallel()
	var _ apieth.ERC20er = apierc20impl.NewERC20(nil, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
}

// TestNewERC20_AcceptsEthereumer verifies that NewERC20 accepts an apieth.Ethereumer
// as its first parameter. apieth.Ethereumer satisfies the ethNodeAPI interface, so
// any Ethereumer is always an acceptable first argument.
func TestNewERC20_AcceptsEthereumer(t *testing.T) {
	t.Parallel()
	var eth apieth.Ethereumer
	result := apierc20impl.NewERC20(eth, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	assert.NotNil(t, result)
}

// ---------------------------------------------------------------------------
// Task 3.3 – CreateRawTransactionEIP1559
// ---------------------------------------------------------------------------

// TestCreateRawTransactionEIP1559_CallsSuggestGasTipCap verifies that
// CreateRawTransactionEIP1559 calls SuggestGasTipCap on the eth field.
func TestCreateRawTransactionEIP1559_CallsSuggestGasTipCap(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		tipCapErr: errors.New("tip cap unavailable"), // cause early return after the call
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	//nolint:dogsled // only the side-effect (call tracking) matters here
	_, _, _ = erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	assert.True(t, stub.suggestGasTipCapCalled, "SuggestGasTipCap must be called")
}

// TestCreateRawTransactionEIP1559_ErrorWhenTipCapFails verifies error propagation
// when SuggestGasTipCap fails.
func TestCreateRawTransactionEIP1559_ErrorWhenTipCapFails(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		tipCapErr: errors.New("tip cap unavailable"),
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	_, _, err := erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	require.Error(t, err)
}

// TestCreateRawTransactionEIP1559_ErrorWhenBlockNumberFails verifies error propagation
// when BlockNumber fails after a successful tip cap call (task 3.3).
func TestCreateRawTransactionEIP1559_ErrorWhenBlockNumberFails(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		tipCap:         big.NewInt(1_000_000_000), // 1 Gwei
		blockNumberErr: errors.New("block number unavailable"),
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	_, _, err := erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	require.Error(t, err)
}

// TestCreateRawTransactionEIP1559_ErrorWhenBaseFeeAbsent verifies that when the
// block header does not contain BaseFeePerGas (pre-London network), an error is returned.
func TestCreateRawTransactionEIP1559_ErrorWhenBaseFeeAbsent(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		tipCap:      big.NewInt(1_000_000_000),
		blockNumber: big.NewInt(100),
		blockInfo:   &dtoeth.BlockInfo{BaseFeePerGas: nil},
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	_, _, err := erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	require.Error(t, err)
}
