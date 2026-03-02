package eth_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
)

// TestEncodeTxDecodeTx_LegacyRoundtrip verifies that a legacy (Type 0) transaction
// can be encoded and decoded without data loss.
func TestEncodeTxDecodeTx_LegacyRoundtrip(t *testing.T) {
	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    42,
		To:       &to,
		Value:    big.NewInt(1_000_000_000_000_000_000), // 1 ETH in wei
		Gas:      21000,
		GasPrice: big.NewInt(20_000_000_000), // 20 Gwei
	})

	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)
	require.NotNil(t, txHex)
	assert.NotEmpty(t, *txHex)

	decoded, err := pkgeth.DecodeTx(*txHex)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	assert.Equal(t, tx.Nonce(), decoded.Nonce())
	assert.Equal(t, tx.To(), decoded.To())
	assert.Equal(t, tx.Value(), decoded.Value())
	assert.Equal(t, tx.Gas(), decoded.Gas())
	assert.Equal(t, tx.GasPrice(), decoded.GasPrice())
	assert.Equal(t, types.LegacyTxType, int(decoded.Type()))
}

// TestEncodeTxDecodeTx_EIP1559Roundtrip verifies that an EIP-1559 (Type 2) transaction
// can be encoded and decoded without data loss using EIP-2718 binary format.
func TestEncodeTxDecodeTx_EIP1559Roundtrip(t *testing.T) {
	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	chainID := big.NewInt(1337)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     7,
		To:        &to,
		Value:     big.NewInt(500_000_000_000_000_000), // 0.5 ETH in wei
		Gas:       21000,
		GasTipCap: big.NewInt(2_000_000_000),  // 2 Gwei tip
		GasFeeCap: big.NewInt(30_000_000_000), // 30 Gwei max fee
	})

	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)
	require.NotNil(t, txHex)
	assert.NotEmpty(t, *txHex)

	decoded, err := pkgeth.DecodeTx(*txHex)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	assert.Equal(t, types.DynamicFeeTxType, int(decoded.Type()))
	assert.Equal(t, tx.Nonce(), decoded.Nonce())
	assert.Equal(t, tx.To(), decoded.To())
	assert.Equal(t, tx.Value(), decoded.Value())
	assert.Equal(t, tx.Gas(), decoded.Gas())
	assert.Equal(t, tx.GasTipCap(), decoded.GasTipCap())
	assert.Equal(t, tx.GasFeeCap(), decoded.GasFeeCap())
	assert.Equal(t, tx.ChainId(), decoded.ChainId())
}

// TestEncodeTx_EIP1559StartsWithTypeByte verifies that an encoded EIP-1559 (Type 2)
// transaction starts with the 0x02 type byte, confirming EIP-2718 binary format.
// With rlp.EncodeToBytes, the encoding would NOT start with 0x02 (it wraps in RLP string).
func TestEncodeTx_EIP1559StartsWithTypeByte(t *testing.T) {
	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(1337),
		Nonce:     1,
		To:        &to,
		Value:     big.NewInt(1),
		Gas:       21000,
		GasTipCap: big.NewInt(2_000_000_000),
		GasFeeCap: big.NewInt(30_000_000_000),
	})

	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)
	require.NotNil(t, txHex)

	// EIP-2718 format: type 2 tx hex must start with "0x02"
	assert.True(t, len(*txHex) >= 4, "encoded tx hex too short")
	assert.Equal(t, "0x02", (*txHex)[:4],
		"EIP-1559 tx must start with 0x02 (EIP-2718 type byte), got: %s", (*txHex)[:4])
}

// TestDecodeTx_InvalidHex verifies that an invalid hex string returns an error.
func TestDecodeTx_InvalidHex(t *testing.T) {
	_, err := pkgeth.DecodeTx("not-valid-hex")
	assert.Error(t, err)
}
