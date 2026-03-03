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

// FromGWei converts GWei(int64) to Wei(*big.Int)
func FromGWei(v int64) *big.Int {
	return big.NewInt(v * params.GWei)
}

// FromFloatEther converts Ether(float64) to Wei(*big.Int)
func FromFloatEther(v float64) *big.Int {
	return big.NewInt(int64(v * params.Ether))
}
