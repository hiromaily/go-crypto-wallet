package eth

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
)

// DecodeBig to handle different response of API between Geth and Parity
func (*Ethereum) DecodeBig(input string) (*big.Int, error) {
	if input == "" || input == "0x" {
		input = "0x0"
	}
	return hexutil.DecodeBig(input)
}

// ValidateAddr validates address
func (*Ethereum) ValidateAddr(addr string) error {
	// validation check
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("address:%s is invalid", addr)
	}
	return nil
}

// Wei   = 1
// GWei  = 1e9  (Giga)
// Ether = 1e18

// Wei :  1000000000000000000
// GWei:  1000000000
// Ether: 1

// FromWei converts Wei(int64) to Wei(*big.Int)
func (*Ethereum) FromWei(v int64) *big.Int {
	return big.NewInt(v * params.Wei)
}

// FromGWei converts GWei(int64) to Wei(*big.Int)
func (*Ethereum) FromGWei(v int64) *big.Int {
	return big.NewInt(v * params.GWei)
}

// FromEther converts Ether(int64) to Wei(*big.Int)
// func (e *Ethereum) FromEther(v int64) *big.Int {
//	return big.NewInt(v * params.Ether)
//}

// FromFloatEther converts Ether(float64) to Wei(*big.Int)
func (*Ethereum) FromFloatEther(v float64) *big.Int {
	return big.NewInt(int64(v * params.Ether))
}

// FloatToBigInt is alias of FromFloatEther for interface
func (e *Ethereum) FloatToBigInt(v float64) *big.Int {
	return e.FromFloatEther(v)
}
