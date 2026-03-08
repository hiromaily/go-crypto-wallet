package eth_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	keygeneth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
)

// TestComputeSafeTxHash_TypehashConstants verifies the pre-computed typehash constants match
// the known values from Safe v1.4.1 contract source code.
// Reference: https://github.com/safe-global/safe-smart-account/blob/v1.4.1/contracts/Safe.sol
func TestComputeSafeTxHash_TypehashConstants(t *testing.T) {
	t.Parallel()

	// Known values from Safe v1.4.1 contract:
	// keccak256("EIP712Domain(uint256 chainId,address verifyingContract)")
	wantDomainTypehash := "47e79534a245952e8b16893a336b85a3d9ea9fa8c573f3d803afb92a79469218"
	// Note: `address` (not `address payable`) — ABI canonical type in getTransactionHash.
	wantSafeTxTypehash := "bb8310d486368db6bd6f849402fdd73ad53d316b5a4b2644ad6efe0f941286d8"

	gotDomain := hex.EncodeToString(crypto.Keccak256([]byte("EIP712Domain(uint256 chainId,address verifyingContract)")))
	safeTxTypeStr := "SafeTx(address to,uint256 value,bytes data,uint8 operation," +
		"uint256 safeTxGas,uint256 baseGas,uint256 gasPrice," +
		"address gasToken,address refundReceiver,uint256 nonce)"
	gotSafeTx := hex.EncodeToString(crypto.Keccak256([]byte(safeTxTypeStr)))

	assert.Equal(t, wantDomainTypehash, gotDomain, "DOMAIN_SEPARATOR_TYPEHASH must match Safe v1.4.1")
	assert.Equal(t, wantSafeTxTypehash, gotSafeTx, "SAFE_TX_TYPEHASH must match Safe v1.4.1")
}

// TestComputeSafeTxHash_KnownVector verifies the full hash computation against a
// reference value computed step-by-step using the canonical EIP-712 algorithm.
//
// This vector uses Anvil defaults (chain ID 31337) so it can be cross-validated
// against a deployed Safe contract in integration tests.
func TestComputeSafeTxHash_KnownVector(t *testing.T) {
	t.Parallel()

	// Test parameters (matching values used in send_multisig_transaction_test.go)
	const (
		chainID     = 31337
		safeAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
		toAddress   = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
		valueStr    = "100000000000000000" // 0.1 ETH in Wei
		nonceStr    = "3"
		zeroAddr    = "0x0000000000000000000000000000000000000000"
	)

	f := &dtoeth.ETHMultisigTransactionFile{
		Version:        1,
		TxType:         "unsigned",
		UUID:           "550e8400-e29b-41d4-a716-446655440000",
		SafeAddress:    safeAddress,
		To:             toAddress,
		Value:          valueStr,
		Data:           "0x",
		Operation:      0,
		SafeTxGas:      "0",
		BaseGas:        "0",
		GasPrice:       "0",
		GasToken:       zeroAddr,
		RefundReceiver: zeroAddr,
		Nonce:          nonceStr,
		ChainID:        chainID,
		SafeTxHash:     "0x" + strings.Repeat("ab", 32), // placeholder; not used in hash computation
		Threshold:      2,
		Signatures:     []dtoeth.SignEntry{},
	}

	got, err := keygeneth.ComputeSafeTxHash(f)
	require.NoError(t, err)

	// Compute expected hash using the reference algorithm (step-by-step).
	expected := computeExpectedSafeTxHash(t, chainID, safeAddress, toAddress, valueStr, nonceStr)

	assert.Equal(t, expected, got, "computed safeTxHash must match reference computation")
}

// computeExpectedSafeTxHash is a step-by-step reference implementation of the EIP-712 algorithm
// used as the expected value in tests. It makes each step explicit to serve as documentation.
func computeExpectedSafeTxHash(
	t *testing.T,
	chainID uint64,
	safeAddress, toAddress, valueStr, nonceStr string,
) [32]byte {
	t.Helper()

	domainTypehash := crypto.Keccak256([]byte(
		"EIP712Domain(uint256 chainId,address verifyingContract)",
	))
	safeTxTypehash := crypto.Keccak256([]byte(
		"SafeTx(address to,uint256 value,bytes data,uint8 operation," +
			"uint256 safeTxGas,uint256 baseGas,uint256 gasPrice," +
			"address gasToken,address refundReceiver,uint256 nonce)",
	))

	// 1. domain separator = keccak256(DOMAIN_SEPARATOR_TYPEHASH || chainId || safeAddr)
	domainInput := make([]byte, 96)
	copy(domainInput[0:32], domainTypehash)
	chainIDPad := make([]byte, 32)
	chainBig := new(big.Int).SetUint64(chainID)
	b := chainBig.Bytes()
	copy(chainIDPad[32-len(b):], b)
	copy(domainInput[32:64], chainIDPad)
	safeAddr := common.HexToAddress(safeAddress)
	copy(domainInput[76:96], safeAddr.Bytes()) // 12 zero bytes + 20 addr bytes
	domainSeparator := crypto.Keccak256(domainInput)

	// 2. struct hash = keccak256(SAFE_TX_TYPEHASH || to || value || keccak256(data) ||
	//                             operation || safeTxGas || baseGas || gasPrice ||
	//                             gasToken || refundReceiver || nonce)
	value, ok := new(big.Int).SetString(valueStr, 10)
	require.True(t, ok)
	nonce, ok := new(big.Int).SetString(nonceStr, 10)
	require.True(t, ok)
	zero := new(big.Int)
	zeroAddr := common.Address{}

	structInput := make([]byte, 352) // 11 × 32
	copy(structInput[0:32], safeTxTypehash)
	toAddr := common.HexToAddress(toAddress)
	copy(structInput[44:64], toAddr.Bytes()) // 12 zero bytes + 20 addr bytes
	vb := value.Bytes()
	copy(structInput[96-len(vb):96], vb)
	// keccak256([]byte{}) = empty data
	emptyDataHash := crypto.Keccak256([]byte{})
	copy(structInput[96:128], emptyDataHash)
	// operation = 0 → zero slot
	// safeTxGas = 0, baseGas = 0, gasPrice = 0 → zero slots
	_ = zero
	_ = zeroAddr
	// gasToken = zero address → zero slot at [256:288]
	// refundReceiver = zero address → zero slot at [288:320]
	nb := nonce.Bytes()
	copy(structInput[352-len(nb):352], nb)
	structHash := crypto.Keccak256(structInput)

	// 3. safeTxHash = keccak256(0x19 0x01 || domainSeparator || structHash)
	final := make([]byte, 66)
	final[0] = 0x19
	final[1] = 0x01
	copy(final[2:34], domainSeparator)
	copy(final[34:66], structHash)

	var result [32]byte
	copy(result[:], crypto.Keccak256(final))
	return result
}

// TestComputeSafeTxHash_Deterministic verifies that the same input always produces the same output.
func TestComputeSafeTxHash_Deterministic(t *testing.T) {
	t.Parallel()

	f := validEIP712File()

	h1, err := keygeneth.ComputeSafeTxHash(f)
	require.NoError(t, err)

	h2, err := keygeneth.ComputeSafeTxHash(f)
	require.NoError(t, err)

	assert.Equal(t, h1, h2, "same input must produce same hash")
}

// TestComputeSafeTxHash_DifferentNonce verifies that different nonces produce different hashes.
func TestComputeSafeTxHash_DifferentNonce(t *testing.T) {
	t.Parallel()

	f0 := validEIP712File()
	f0.Nonce = "0"

	f1 := validEIP712File()
	f1.Nonce = "1"

	h0, err := keygeneth.ComputeSafeTxHash(f0)
	require.NoError(t, err)

	h1, err := keygeneth.ComputeSafeTxHash(f1)
	require.NoError(t, err)

	assert.NotEqual(t, h0, h1, "different nonce must produce different hash")
}

// TestComputeSafeTxHash_DifferentChainID verifies that different chain IDs produce different hashes.
func TestComputeSafeTxHash_DifferentChainID(t *testing.T) {
	t.Parallel()

	f1 := validEIP712File()
	f1.ChainID = 1

	f137 := validEIP712File()
	f137.ChainID = 137 // Polygon

	h1, err := keygeneth.ComputeSafeTxHash(f1)
	require.NoError(t, err)

	h137, err := keygeneth.ComputeSafeTxHash(f137)
	require.NoError(t, err)

	assert.NotEqual(t, h1, h137, "different chainIDs must produce different hashes")
}

// TestComputeSafeTxHash_InvalidValue verifies that an unparseable Value returns an error.
func TestComputeSafeTxHash_InvalidValue(t *testing.T) {
	t.Parallel()

	f := validEIP712File()
	f.Value = "not-a-number"

	_, err := keygeneth.ComputeSafeTxHash(f)
	assert.Error(t, err)
}

// TestComputeSafeTxHash_InvalidNonce verifies that an unparseable Nonce returns an error.
func TestComputeSafeTxHash_InvalidNonce(t *testing.T) {
	t.Parallel()

	f := validEIP712File()
	f.Nonce = "not-a-number"

	_, err := keygeneth.ComputeSafeTxHash(f)
	assert.Error(t, err)
}

// TestComputeSafeTxHash_InvalidData verifies that an invalid hex Data returns an error.
func TestComputeSafeTxHash_InvalidData(t *testing.T) {
	t.Parallel()

	f := validEIP712File()
	f.Data = "0xZZZZ" // invalid hex

	_, err := keygeneth.ComputeSafeTxHash(f)
	assert.Error(t, err)
}

// TestComputeSafeTxHash_WithData verifies that non-empty data is correctly handled
// (keccak256(data) replaces the data in the struct hash).
func TestComputeSafeTxHash_WithData(t *testing.T) {
	t.Parallel()

	fEmpty := validEIP712File()
	fEmpty.Data = "0x"

	fWithData := validEIP712File()
	fWithData.Data = "0xdeadbeef"

	hEmpty, err := keygeneth.ComputeSafeTxHash(fEmpty)
	require.NoError(t, err)

	hData, err := keygeneth.ComputeSafeTxHash(fWithData)
	require.NoError(t, err)

	assert.NotEqual(t, hEmpty, hData, "data affects the hash")
}

// validEIP712File returns a minimal valid ETHMultisigTransactionFile for EIP-712 tests.
func validEIP712File() *dtoeth.ETHMultisigTransactionFile {
	return &dtoeth.ETHMultisigTransactionFile{
		Version:        1,
		TxType:         "unsigned",
		UUID:           "550e8400-e29b-41d4-a716-446655440000",
		SafeAddress:    "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		To:             "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
		Value:          "100000000000000000",
		Data:           "0x",
		Operation:      0,
		SafeTxGas:      "0",
		BaseGas:        "0",
		GasPrice:       "0",
		GasToken:       "0x0000000000000000000000000000000000000000",
		RefundReceiver: "0x0000000000000000000000000000000000000000",
		Nonce:          "3",
		ChainID:        31337,
		SafeTxHash:     "0x" + strings.Repeat("ab", 32),
		Threshold:      2,
		Signatures:     []dtoeth.SignEntry{},
	}
}
