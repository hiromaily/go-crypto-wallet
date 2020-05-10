package eth

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
)

func castToInt64(val interface{}) (int64, error) {
	v, err := castToString(val)
	if err != nil {
		return 0, err
	}

	v2, err := hexutil.DecodeBig(v)
	if err != nil {
		return 0, err
	}
	return v2.Int64(), nil
}

func castToString(val interface{}) (string, error) {
	v, ok := val.(string)
	if !ok {
		return "", errors.New("fail to cast to string")
	}
	return v, nil
}

func castToSliceString(val interface{}) ([]string, error) {
	data, ok := val.([]interface{})
	if !ok {
		return nil, errors.New("fail to cast to []string")
	}
	if len(data) == 0 {
		return []string{}, nil
	}

	var ret []string
	for _, v := range data {
		sRet, err := castToString(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, sRet)
	}
	return ret, nil
}

// DecodeBig to handle different response of API between Geth and Parity
func (e *Ethereum) DecodeBig(input string) (*big.Int, error) {
	if input == "" || input == "0x" {
		input = "0x0"
	}
	if e.isParity {
		i, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return nil, err
		}
		return big.NewInt(i), nil
	}
	return hexutil.DecodeBig(input)
}

func setZeroHex(input string) string {
	if input == "" || input == "0x" {
		input = "0x0"
	}
	return input
}

func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

// ValidationAddr validates address
func (e *Ethereum) ValidationAddr(addr string) error {
	// validation check
	if !common.IsHexAddress(addr) {
		return errors.Errorf("address:%s is invalid", addr)
	}
	return nil
}

//Wei   = 1
//GWei  = 1e9  (Giga)
//Ether = 1e18

//Wei :  1000000000000000000
//GWei:  1000000000
//Ether: 1

// FromWei converts Wei(int64) to Wei(*big.Int)
func (e *Ethereum) FromWei(v int64) *big.Int {
	return big.NewInt(v * params.Wei)
}

// FromGWei converts GWei(int64) to Wei(*big.Int)
func (e *Ethereum) FromGWei(v int64) *big.Int {
	return big.NewInt(v * params.GWei)
}

// FromEther converts Ether(int64) to Wei(*big.Int)
func (e *Ethereum) FromEther(v int64) *big.Int {
	return big.NewInt(v * params.Ether)
}
