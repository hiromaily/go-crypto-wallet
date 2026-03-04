package erc20_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	apierc20impl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/erc20"
	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// stubEthNode is a minimal test double for the eth node provider (ethNodeAPI).
// It records method calls for verification in TDD tests.
type stubEthNode struct {
	supportsEIP1559 bool
	tipCap          *big.Int
	tipCapErr       error
	blockNumber     *big.Int
	blockNumberErr  error
	blockInfo       *ethrpc.BlockInfo
	blockInfoErr    error
	txCount         *big.Int
	txCountErr      error

	// call trackers
	suggestGasTipCapCalled bool
	blockNumberCalled      bool
	getBlockByNumberCalled bool
}

func (s *stubEthNode) SupportsEIP1559(_ context.Context) bool { return s.supportsEIP1559 }

func (s *stubEthNode) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	s.suggestGasTipCapCalled = true
	return s.tipCap, s.tipCapErr
}

func (s *stubEthNode) BlockNumber(_ context.Context) (*big.Int, error) {
	s.blockNumberCalled = true
	return s.blockNumber, s.blockNumberErr
}

func (s *stubEthNode) GetBlockByNumber(_ context.Context, _ uint64) (*ethrpc.BlockInfo, error) {
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
// Task 3.2 – SupportsEIP1559 delegation
// ---------------------------------------------------------------------------

// TestERC20_SupportsEIP1559 verifies that ERC20.SupportsEIP1559 delegates to the
// eth field, returning its result directly (task 3.2).
func TestERC20_SupportsEIP1559(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		supportsEIP1559 bool
		want            bool
	}{
		{"returns true when eth supports EIP-1559", true, true},
		{"returns false when eth does not support EIP-1559", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stub := &stubEthNode{supportsEIP1559: tt.supportsEIP1559}
			erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
			got := erc20.SupportsEIP1559(context.Background())
			assert.Equal(t, tt.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// Task 3.3 – CreateRawTransactionEIP1559 EIP-1559 path
// ---------------------------------------------------------------------------

// TestCreateRawTransactionEIP1559_CallsSuggestGasTipCapWhenEIP1559Supported verifies
// that when the node supports EIP-1559, the method calls SuggestGasTipCap on the
// eth field (task 3.3 RED: current stub delegation never calls SuggestGasTipCap).
func TestCreateRawTransactionEIP1559_CallsSuggestGasTipCapWhenEIP1559Supported(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		supportsEIP1559: true,
		tipCapErr:       errors.New("tip cap unavailable"), // cause early return after the call
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	//nolint:dogsled // only the side-effect (call tracking) matters here
	_, _, _ = erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	assert.True(t, stub.suggestGasTipCapCalled, "SuggestGasTipCap must be called when EIP-1559 is supported")
}

// TestCreateRawTransactionEIP1559_DoesNotCallSuggestGasTipCapWhenEIP1559NotSupported
// verifies that the fallback path (EIP-1559 not supported) does NOT call
// SuggestGasTipCap — it should fall straight through to CreateRawTransaction.
func TestCreateRawTransactionEIP1559_DoesNotCallSuggestGasTipCapWhenEIP1559NotSupported(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{supportsEIP1559: false}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	//nolint:dogsled // only the side-effect (call tracking) matters here
	_, _, _ = erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	assert.False(t, stub.suggestGasTipCapCalled, "SuggestGasTipCap must NOT be called on the legacy fallback path")
}

// TestCreateRawTransactionEIP1559_ErrorWhenTipCapFails verifies error propagation
// when SuggestGasTipCap fails (task 3.3).
func TestCreateRawTransactionEIP1559_ErrorWhenTipCapFails(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		supportsEIP1559: true,
		tipCapErr:       errors.New("tip cap unavailable"),
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
		supportsEIP1559: true,
		tipCap:          big.NewInt(1_000_000_000), // 1 Gwei
		blockNumberErr:  errors.New("block number unavailable"),
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	_, _, err := erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	require.Error(t, err)
}

// TestCreateRawTransactionEIP1559_ErrorWhenBaseFeeAbsent verifies that when the
// block header does not contain BaseFeePerGas (pre-London network), an error is
// returned even though SupportsEIP1559 returned true (task 3.3).
func TestCreateRawTransactionEIP1559_ErrorWhenBaseFeeAbsent(t *testing.T) {
	t.Parallel()
	stub := &stubEthNode{
		supportsEIP1559: true,
		tipCap:          big.NewInt(1_000_000_000),
		blockNumber:     big.NewInt(100),
		blockInfo:       &ethrpc.BlockInfo{BaseFeePerGas: nil},
	}
	erc20 := apierc20impl.NewERC20(stub, nil, nil, domainCoin.TokenHYT, nil, "", "", "", 0)
	_, _, err := erc20.CreateRawTransactionEIP1559(context.Background(), "", "", 0, 0)
	require.Error(t, err)
}
