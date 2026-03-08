package eth

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
)

// domainSeparatorTypehash is the keccak256 hash of the EIP-712 domain type string.
// Matches the constant in Safe v1.4.1 contract:
//
//	keccak256("EIP712Domain(uint256 chainId,address verifyingContract)")
//	= 0x47e79534a245952e8b16893a336b85a3d9ea9fa8c573f3d803afb92a79469218
var domainSeparatorTypehash = mustHash("EIP712Domain(uint256 chainId,address verifyingContract)")

// safeTxTypehash is the keccak256 hash of the Safe transaction struct type string.
// Matches the constant in Safe v1.4.1 contract
// = 0xbb8310d486368db6bd6f849402fdd73ad53d316b5a4b2644ad6efe0f941286d8.
//
// Note: `refundReceiver` uses the ABI canonical type `address` (not `address payable`),
// which is how the Safe contract encodes it in getTransactionHash and encodeTransactionData.
var safeTxTypehash = mustHash(
	"SafeTx(address to,uint256 value,bytes data,uint8 operation," +
		"uint256 safeTxGas,uint256 baseGas,uint256 gasPrice," +
		"address gasToken,address refundReceiver,uint256 nonce)",
)

func mustHash(s string) [32]byte {
	var b [32]byte
	copy(b[:], crypto.Keccak256([]byte(s)))
	return b
}

// ComputeSafeTxHash recomputes the EIP-712 safeTxHash from the fields stored in
// an ETHMultisigTransactionFile. This is a pure-Go function that makes no network
// calls — suitable for offline signing on air-gapped Keygen and Sign wallets.
//
// Algorithm (Safe v1.4.1 EIP-712):
//
//  1. domainSeparator = keccak256(DOMAIN_SEPARATOR_TYPEHASH || uint256(chainId) || address(safeAddr))
//  2. structHash      = keccak256(SAFE_TX_TYPEHASH || address(to) || uint256(value) ||
//     keccak256(data) || uint256(operation) || uint256(safeTxGas) ||
//     uint256(baseGas) || uint256(gasPrice) || address(gasToken) ||
//     address(refundReceiver) || uint256(nonce))
//  3. safeTxHash      = keccak256(0x19 0x01 || domainSeparator || structHash)
//
// The caller should compare the returned hash to f.SafeTxHash and abort with
// dtoeth.ErrSafeTxHashMismatch if they differ.
func ComputeSafeTxHash(f *dtoeth.ETHMultisigTransactionFile) ([32]byte, error) {
	value, ok := new(big.Int).SetString(f.Value, 10)
	if !ok {
		return [32]byte{}, fmt.Errorf("invalid value field %q", f.Value)
	}
	nonce, ok := new(big.Int).SetString(f.Nonce, 10)
	if !ok {
		return [32]byte{}, fmt.Errorf("invalid nonce field %q", f.Nonce)
	}
	safeTxGas, ok := new(big.Int).SetString(f.SafeTxGas, 10)
	if !ok {
		return [32]byte{}, fmt.Errorf("invalid safe_tx_gas field %q", f.SafeTxGas)
	}
	baseGas, ok := new(big.Int).SetString(f.BaseGas, 10)
	if !ok {
		return [32]byte{}, fmt.Errorf("invalid base_gas field %q", f.BaseGas)
	}
	gasPrice, ok := new(big.Int).SetString(f.GasPrice, 10)
	if !ok {
		return [32]byte{}, fmt.Errorf("invalid gas_price field %q", f.GasPrice)
	}

	data, err := decodeHexBytesEIP712(f.Data)
	if err != nil {
		return [32]byte{}, fmt.Errorf("invalid data field: %w", err)
	}

	chainID := new(big.Int).SetUint64(f.ChainID)
	safeAddr := common.HexToAddress(f.SafeAddress)
	toAddr := common.HexToAddress(f.To)
	gasTokenAddr := common.HexToAddress(f.GasToken)
	refundAddr := common.HexToAddress(f.RefundReceiver)

	domainSep := computeDomainSeparator(chainID, safeAddr)
	structHash := computeStructHash(
		toAddr, value, data, f.Operation,
		safeTxGas, baseGas, gasPrice,
		gasTokenAddr, refundAddr, nonce,
	)

	return computeEIP191Hash(domainSep, structHash), nil
}

// computeDomainSeparator computes the EIP-712 domain separator:
//
//	keccak256(DOMAIN_SEPARATOR_TYPEHASH || uint256(chainId) || address(safeAddr))
func computeDomainSeparator(chainID *big.Int, safeAddr common.Address) [32]byte {
	// 3 × 32 = 96 bytes
	packed := make([]byte, 96)
	copy(packed[0:32], domainSeparatorTypehash[:])
	copy(packed[32:64], padUint256(chainID))
	copy(packed[64:96], padAddress(safeAddr))

	var result [32]byte
	copy(result[:], crypto.Keccak256(packed))
	return result
}

// computeStructHash computes the EIP-712 struct hash for a Safe transaction:
//
//	keccak256(SAFE_TX_TYPEHASH || address(to) || uint256(value) || keccak256(data) ||
//	           uint256(operation) || uint256(safeTxGas) || uint256(baseGas) ||
//	           uint256(gasPrice) || address(gasToken) || address(refundReceiver) ||
//	           uint256(nonce))
//
// Note: dynamic type `bytes data` is replaced by keccak256(data) per the EIP-712 spec.
func computeStructHash(
	to common.Address,
	value *big.Int,
	data []byte,
	operation uint8,
	safeTxGas, baseGas, gasPrice *big.Int,
	gasToken, refundReceiver common.Address,
	nonce *big.Int,
) [32]byte {
	// 11 fields × 32 bytes = 352 bytes
	packed := make([]byte, 352)
	copy(packed[0:32], safeTxTypehash[:])
	copy(packed[32:64], padAddress(to))
	copy(packed[64:96], padUint256(value))
	copy(packed[96:128], crypto.Keccak256(data)) // bytes → keccak256(data)
	copy(packed[128:160], padUint256(new(big.Int).SetUint64(uint64(operation))))
	copy(packed[160:192], padUint256(safeTxGas))
	copy(packed[192:224], padUint256(baseGas))
	copy(packed[224:256], padUint256(gasPrice))
	copy(packed[256:288], padAddress(gasToken))
	copy(packed[288:320], padAddress(refundReceiver))
	copy(packed[320:352], padUint256(nonce))

	var result [32]byte
	copy(result[:], crypto.Keccak256(packed))
	return result
}

// computeEIP191Hash computes the final EIP-712 hash:
//
//	keccak256(0x19 0x01 || domainSeparator || structHash)
func computeEIP191Hash(domainSeparator, structHash [32]byte) [32]byte {
	packed := make([]byte, 66) // 2 + 32 + 32
	packed[0] = 0x19
	packed[1] = 0x01
	copy(packed[2:34], domainSeparator[:])
	copy(packed[34:66], structHash[:])

	var result [32]byte
	copy(result[:], crypto.Keccak256(packed))
	return result
}

// padAddress right-aligns a 20-byte Ethereum address into a 32-byte ABI slot
// (12 zero bytes of padding on the left).
func padAddress(addr common.Address) []byte {
	padded := make([]byte, 32)
	copy(padded[12:], addr.Bytes())
	return padded
}

// padUint256 big-endian encodes n into a 32-byte ABI slot (left-zero-padded).
func padUint256(n *big.Int) []byte {
	padded := make([]byte, 32)
	b := n.Bytes()
	if len(b) > 32 {
		b = b[len(b)-32:]
	}
	copy(padded[32-len(b):], b)
	return padded
}

// decodeHexBytesEIP712 decodes a 0x-prefixed hex string to bytes.
// An empty "0x" returns an empty byte slice (used for simple ETH transfers).
func decodeHexBytesEIP712(hexStr string) ([]byte, error) {
	stripped := strings.TrimPrefix(hexStr, "0x")
	if stripped == "" {
		return []byte{}, nil
	}
	return hex.DecodeString(stripped)
}
