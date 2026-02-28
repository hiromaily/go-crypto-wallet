// Package erc20 — internal unit tests for createTransferData.
// Uses the same package so the private helper can be called directly.
package erc20

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// transferMethodSelector is the ABI method selector for transfer(address,uint256).
// keccak256("transfer(address,uint256)")[0:4] == 0xa9059cbb
var transferMethodSelector = []byte{0xa9, 0x05, 0x9c, 0xbb}

// toAddr is a valid, non-zero Ethereum address used across transfer-data tests.
const toAddrHex = "0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF"

// TestCreateTransferData_HasMethodSelector verifies that the first 4 bytes of
// the encoded calldata are the transfer(address,uint256) method selector 0xa9059cbb.
// This is the key property validated by task 5.2.
func TestCreateTransferData_HasMethodSelector(t *testing.T) {
	t.Parallel()
	e := &ERC20{}
	data := e.createTransferData(toAddrHex, big.NewInt(100))
	require.GreaterOrEqual(t, len(data), 4, "calldata must contain at least the 4-byte method selector")
	assert.Equal(t, transferMethodSelector, data[:4],
		"first 4 bytes must be the transfer() ABI method selector 0xa9059cbb")
}

// TestCreateTransferData_Length verifies the ABI-encoded calldata is exactly
// 68 bytes: 4-byte selector + 32-byte padded address + 32-byte amount.
func TestCreateTransferData_Length(t *testing.T) {
	t.Parallel()
	e := &ERC20{}
	data := e.createTransferData(toAddrHex, big.NewInt(1_000_000))
	assert.Equal(t, 68, len(data), "ABI-encoded transfer calldata must be 68 bytes")
}

// TestCreateTransferData_ZeroAmount verifies zero-amount transfers also produce
// the correct method selector (used when sending the full balance).
func TestCreateTransferData_ZeroAmount(t *testing.T) {
	t.Parallel()
	e := &ERC20{}
	data := e.createTransferData(toAddrHex, big.NewInt(0))
	require.GreaterOrEqual(t, len(data), 4)
	assert.Equal(t, transferMethodSelector, data[:4])
}
