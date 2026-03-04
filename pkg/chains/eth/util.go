package eth

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

// DecodeBig decodes a hex-encoded big.Int, treating empty or "0x" as zero.
func DecodeBig(input string) (*big.Int, error) {
	if input == "" || input == "0x" {
		input = "0x0"
	}
	return hexutil.DecodeBig(input)
}

// ValidateAddr validates an Ethereum hex address.
func ValidateAddr(addr string) error {
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("address:%s is invalid", addr)
	}
	return nil
}

// Wei   = 1
// GWei  = 1e9  (Giga)
// Ether = 1e18

// FromWei converts Wei(int64) to Wei(*big.Int)
func FromWei(v int64) *big.Int {
	return big.NewInt(v * params.Wei)
}

// FromGWei converts GWei(int64) to Wei(*big.Int).
// Uses big.Int multiplication to avoid int64 overflow for large values.
func FromGWei(v int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(v), big.NewInt(params.GWei))
}

// FromFloatEther converts Ether(float64) to Wei(*big.Int).
// Uses big.Float with float64 precision (53 bits) to match float64 rounding
// while converting directly to big.Int to avoid int64 overflow.
func FromFloatEther(v float64) *big.Int {
	bf := new(big.Float).SetPrec(53).SetFloat64(v)
	ether := new(big.Float).SetPrec(53).SetFloat64(params.Ether)
	result, _ := new(big.Float).SetPrec(53).Mul(bf, ether).Int(nil)
	return result
}
