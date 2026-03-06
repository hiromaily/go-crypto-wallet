package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GetAddressesByLabel returns addresses of account(label)
// Note: even if client has 5 addresses, it returns 15 addresses
//
//	it seems 3 different address types are returned respectively
//
// For now, it would be better to stop using it
func (b *Bitcoin) GetAddressesByLabel(labelName string) ([]btcutil.Address, error) {
	logger.Debug("getting addresses by label", "label", labelName)

	labels, err := b.pkgrpc.GetAddressesByLabel(labelName)
	if err != nil {
		logger.Debug("getaddressesbylabel RPC failed", "label", labelName, "error", err)
		return nil, fmt.Errorf("fail to call btcrpc.GetAddressesByLabel(%s): %w", labelName, err)
	}

	if len(labels) == 0 {
		logger.Debug("no addresses found for label", "label", labelName)
		return nil, nil
	}

	logger.Debug("found addresses for label",
		"label", labelName,
		"raw_count", len(labels))

	// retrieve
	resAddrs := make([]btcutil.Address, 0, len(labels))
	for key := range labels {
		// key is address string
		address, err := b.DecodeAddress(key)
		if err != nil {
			logger.Error(
				"fail to call b.DecodeAddress()",
				"address", key,
				"label", labelName,
				"error", err)
			continue
		}

		resAddrs = append(resAddrs, address)
	}

	logger.Debug("successfully decoded addresses",
		"label", labelName,
		"decoded_count", len(resAddrs))

	return resAddrs, nil
}

// ValidateAddress validate address
func (b *Bitcoin) ValidateAddress(addr string) (*dtobtc.ValidateAddressResult, error) {
	result, err := b.pkgrpc.ValidateAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.ValidateAddress(%s): %w", addr, err)
	}
	if !result.IsValid {
		return nil, fmt.Errorf("this address is invalid: %v", result)
	}

	return toValidateAddressResult(result), nil
}

// GetAddressInfo returns information about the given bitcoin address
func (b *Bitcoin) GetAddressInfo(addr string) (*dtobtc.AddressInfo, error) {
	result, err := b.pkgrpc.GetAddressInfo(addr)
	if err != nil {
		return nil, err
	}
	return toAddressInfo(result), nil
}

func toAddressInfo(r *btcrpc.GetAddressInfoResult) *dtobtc.AddressInfo {
	return &dtobtc.AddressInfo{
		Address:             r.Address,
		ScriptPubKey:        r.ScriptPubKey,
		IsMine:              r.IsMine,
		Solvable:            r.Solvable,
		Desc:                r.Desc,
		IsWatchOnly:         r.IsWatchOnly,
		IsScript:            r.IsScript,
		IsWitness:           r.IsWitness,
		Pubkey:              r.Pubkey,
		IsCompressed:        r.IsCompressed,
		IsChange:            r.IsChange,
		HDKeyPath:           r.HDKeyPath,
		HDMasterFingerprint: r.HDMasterFingerprint,
		Labels:              []string(r.Labels),
		Label:               r.Label,
	}
}

func toValidateAddressResult(r *btcrpc.ValidateAddressResult) *dtobtc.ValidateAddressResult {
	return &dtobtc.ValidateAddressResult{
		IsValid:           r.IsValid,
		Address:           r.Address,
		ScriptPubKey:      r.ScriptPubKey,
		IsScript:          r.IsScript,
		IsWitness:         r.IsWitness,
		WitnessVersion:    r.WitnessVersion,
		WitnessProgramHex: r.WitnessProgramHex,
	}
}

// DecodeAddress decode string address to type Address
func (b *Bitcoin) DecodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcutil.DecodeAddress(): %w", err)
	}
	return address, nil
}
