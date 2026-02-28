// Package erc20 — internal unit tests for getNonce.
// Uses the same package so private methods can be called directly.
package erc20

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
)

// nonceStub is a minimal apieth.ERC20NodeAPI test double for nonce tests.
// Only GetTransactionCount needs a real implementation; the remaining methods
// are present solely to satisfy the interface and panic if called unexpectedly.
type nonceStub struct {
	txCount    *big.Int
	txCountErr error
	called     bool
}

func (*nonceStub) SupportsEIP1559(_ context.Context) bool { panic("unexpected call") }
func (*nonceStub) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	panic("unexpected call")
}
func (*nonceStub) BlockNumber(_ context.Context) (*big.Int, error) { panic("unexpected call") }
func (*nonceStub) GetBlockByNumber(_ context.Context, _ uint64) (*domainETH.BlockInfo, error) {
	panic("unexpected call")
}

func (n *nonceStub) GetTransactionCount(_ context.Context, _ string, _ domainETH.QuantityTag) (*big.Int, error) {
	n.called = true
	return n.txCount, n.txCountErr
}

// TestGetNonce_DelegatesToEth verifies that getNonce calls e.eth.GetTransactionCount
// (task 3.4 RED: current implementation calls e.client.PendingNonceAt instead).
func TestGetNonce_DelegatesToEth(t *testing.T) {
	t.Parallel()
	stub := &nonceStub{txCount: big.NewInt(42)}
	e := &ERC20{eth: stub}
	nonce, err := e.getNonce(context.Background(), "0x0000000000000000000000000000000000000001", 0)
	require.NoError(t, err)
	assert.True(t, stub.called, "GetTransactionCount must be called by getNonce")
	assert.Equal(t, uint64(42), nonce)
}

// TestGetNonce_AddsAdditionalNonce verifies that additionalNonce is correctly
// added to the pending transaction count.
func TestGetNonce_AddsAdditionalNonce(t *testing.T) {
	t.Parallel()
	stub := &nonceStub{txCount: big.NewInt(10)}
	e := &ERC20{eth: stub}
	nonce, err := e.getNonce(context.Background(), "0x0000000000000000000000000000000000000001", 3)
	require.NoError(t, err)
	assert.Equal(t, uint64(13), nonce) // 10 + 3
}

// TestGetNonce_PropagatesError verifies that errors from GetTransactionCount are
// returned to the caller.
func TestGetNonce_PropagatesError(t *testing.T) {
	t.Parallel()
	stub := &nonceStub{txCountErr: errors.New("rpc error")}
	e := &ERC20{eth: stub}
	_, err := e.getNonce(context.Background(), "0x0000000000000000000000000000000000000001", 0)
	require.Error(t, err)
}
