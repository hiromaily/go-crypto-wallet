package eth_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
)

func newTestRawTx(t *testing.T, tx *types.Transaction, fromHex string) *pkgeth.RawTx {
	t.Helper()
	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)
	toHex := ""
	if tx.To() != nil {
		toHex = tx.To().Hex()
	}
	return &pkgeth.RawTx{
		UUID:  "test-uuid",
		From:  fromHex,
		To:    toHex,
		Value: *tx.Value(),
		Nonce: tx.Nonce(),
		TxHex: *txHex,
	}
}

// TestSignTxOffline_Legacy verifies offline signing of a legacy (Type 0) transaction.
func TestSignTxOffline_Legacy(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	fromAddr := crypto.PubkeyToAddress(privKey.PublicKey)

	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       &to,
		Value:    big.NewInt(1_000_000),
		Gas:      21000,
		GasPrice: big.NewInt(1_000_000_000),
	})
	chainID := big.NewInt(1337)
	rawTx := newTestRawTx(t, tx, fromAddr.Hex())

	signed, err := pkgeth.SignTxOffline(rawTx, privKey, chainID)
	require.NoError(t, err)
	require.NotNil(t, signed)

	// Decode and verify the signed transaction
	decoded, err := pkgeth.DecodeTx(signed.TxHex)
	require.NoError(t, err)

	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, decoded)
	require.NoError(t, err)
	assert.Equal(t, fromAddr, sender, "sender address must match the signing key")
}

// TestSignTxOffline_EIP1559 verifies offline signing of an EIP-1559 (Type 2) transaction.
func TestSignTxOffline_EIP1559(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	fromAddr := crypto.PubkeyToAddress(privKey.PublicKey)

	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	chainID := big.NewInt(1337)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     3,
		To:        &to,
		Value:     big.NewInt(500_000),
		Gas:       21000,
		GasTipCap: big.NewInt(2_000_000_000),
		GasFeeCap: big.NewInt(30_000_000_000),
	})
	rawTx := newTestRawTx(t, tx, fromAddr.Hex())

	signed, err := pkgeth.SignTxOffline(rawTx, privKey, chainID)
	require.NoError(t, err)
	require.NotNil(t, signed)

	// Verify the signed transaction preserves the type and is properly signed
	decoded, err := pkgeth.DecodeTx(signed.TxHex)
	require.NoError(t, err)
	assert.Equal(t, types.DynamicFeeTxType, int(decoded.Type()))

	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, decoded)
	require.NoError(t, err)
	assert.Equal(t, fromAddr, sender, "sender address must match the signing key")
}

// TestSignTxOffline_NilPrivKey verifies that a nil private key returns an error.
func TestSignTxOffline_NilPrivKey(t *testing.T) {
	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	tx := types.NewTx(&types.LegacyTx{
		Nonce: 0, To: &to, Value: big.NewInt(1), Gas: 21000,
		GasPrice: big.NewInt(1_000_000_000),
	})
	rawTx := newTestRawTx(t, tx, "0xdeadbeef")

	_, err := pkgeth.SignTxOffline(rawTx, nil, big.NewInt(1337))
	assert.Error(t, err)
}

// TestSignTxOffline_ContractCreation verifies that signing a contract creation
// transaction (To == nil) does not panic and produces a valid signed transaction.
func TestSignTxOffline_ContractCreation(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	fromAddr := crypto.PubkeyToAddress(privKey.PublicKey)

	// Contract creation: To is nil
	chainID := big.NewInt(1337)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     0,
		To:        nil, // contract creation
		Value:     big.NewInt(0),
		Gas:       100000,
		GasTipCap: big.NewInt(2_000_000_000),
		GasFeeCap: big.NewInt(30_000_000_000),
		Data:      []byte{0x60, 0x60}, // minimal bytecode stub
	})
	rawTx := newTestRawTx(t, tx, fromAddr.Hex())

	signed, err := pkgeth.SignTxOffline(rawTx, privKey, chainID)
	require.NoError(t, err)
	require.NotNil(t, signed)
	assert.Empty(t, signed.To, "To must be empty string for contract creation")

	// Verify signature recovers the correct sender
	decoded, err := pkgeth.DecodeTx(signed.TxHex)
	require.NoError(t, err)
	signer := types.LatestSignerForChainID(chainID)
	sender, err := types.Sender(signer, decoded)
	require.NoError(t, err)
	assert.Equal(t, fromAddr, sender)
}

// TestSignTxOffline_NilRawTx verifies that a nil rawTx returns an error.
func TestSignTxOffline_NilRawTx(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	_, err = pkgeth.SignTxOffline(nil, privKey, big.NewInt(1337))
	assert.Error(t, err)
}

// TestSignTxOffline_NegativeChainID verifies that a negative chainID returns an error.
func TestSignTxOffline_NegativeChainID(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	tx := types.NewTx(&types.LegacyTx{
		Nonce: 0, To: &to, Value: big.NewInt(1), Gas: 21000,
		GasPrice: big.NewInt(1_000_000_000),
	})
	rawTx := newTestRawTx(t, tx, crypto.PubkeyToAddress(privKey.PublicKey).Hex())
	_, err = pkgeth.SignTxOffline(rawTx, privKey, big.NewInt(-1))
	assert.Error(t, err)
}
