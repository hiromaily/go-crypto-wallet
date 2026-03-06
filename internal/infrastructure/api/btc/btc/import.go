package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// ImportPrivKey import privKey to wallet
// - Rescan  *bool `jsonrpcdefault:"true"`
func (b *Bitcoin) ImportPrivKey(privKeyWIF *btcutil.WIF) error {
	err := b.btcdClient.ImportPrivKey(privKeyWIF)
	if err != nil {
		return fmt.Errorf("fail to call client.ImportPrivKey(): %w", err)
	}

	return nil
}

// ImportPrivKeyLabel import privKey with label to wallet
// - Rescan  *bool `jsonrpcdefault:"true"`
func (b *Bitcoin) ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error {
	err := b.btcdClient.ImportPrivKeyLabel(privKeyWIF, label)
	if err != nil {
		return fmt.Errorf("fail to call client.ImportPrivKeyLabel(): %w", err)
	}

	return nil
}

// ImportPrivKeyWithoutReScan import privKey without rescan to wallet
func (b *Bitcoin) ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error {
	err := b.btcdClient.ImportPrivKeyRescan(privKeyWIF, label, false)
	if err != nil {
		return fmt.Errorf("fail to call ImportPrivKeyRescan(): %w", err)
	}

	return nil
}

// ImportAddress import pubkey to wallet
// Note: This is a legacy wallet method. For descriptor wallets (Bitcoin Core v23.0+),
// consider using importdescriptors instead for better functionality and future compatibility.
func (b *Bitcoin) ImportAddress(pubkey string) error {
	err := b.btcdClient.ImportAddress(pubkey)
	if err != nil {
		return fmt.Errorf("fail to call ImportAddress(): %w", err)
	}

	return nil
}

// ImportAddressWithoutReScan import pubkey without rescan
func (b *Bitcoin) ImportAddressWithoutReScan(pubkey string) error {
	if err := b.ImportAddressWithLabel(pubkey, "", false); err != nil {
		return fmt.Errorf("fail to call ImportAddressWithoutReScan(): %w", err)
	}

	return nil
}

// ImportAddressWithLabel import given address with label to wallet
// - rescan is adjustable
// Note: This is a legacy wallet method. For descriptor wallets (Bitcoin Core v23.0+),
// consider using importdescriptors instead for better functionality and future compatibility.
func (b *Bitcoin) ImportAddressWithLabel(address, label string, rescan bool) error {
	if err := b.pkgrpc.ImportAddress(address, label, rescan); err != nil {
		return fmt.Errorf("fail to call btcrpc.ImportAddress(): %w", err)
	}

	return nil
}

// ImportDescriptors imports output descriptors into the wallet (Bitcoin Core v0.21.0+).
//
// This is the modern replacement for importaddress and importmulti.
// It allows importing descriptors with full scriptPubKey information, enabling
// Bitcoin Core to mark addresses as solvable for transaction creation.
//
// Parameters:
//   - requests: List of descriptor import requests
//
// Returns:
//   - List of responses (one per request, in the same order)
//   - Error if the RPC call fails
//
// Notes:
//   - All descriptors must include checksums
//   - Use "now" timestamp to skip rescanning (fastest)
//   - Set active=true to enable spending from these descriptors
//   - Set watchonly=true for watch-only wallets (no private keys)
//
// Reference: Bitcoin Core RPC documentation - importdescriptors
func (b *Bitcoin) ImportDescriptors(
	requests []dtobtc.ImportDescriptorsRequest,
) ([]dtobtc.ImportDescriptorsResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no descriptors to import")
	}

	wireReqs := make([]btcrpc.ImportDescriptorsRequest, len(requests))
	for i, r := range requests {
		wireReqs[i] = btcrpc.ImportDescriptorsRequest{
			Descriptor: r.Descriptor,
			Timestamp:  r.Timestamp,
			Active:     r.Active,
			Range:      r.Range,
			Label:      r.Label,
			Internal:   r.Internal,
			Watchonly:  r.Watchonly,
		}
	}

	wireResps, err := b.pkgrpc.ImportDescriptors(wireReqs)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.ImportDescriptors(): %w", err)
	}

	if len(wireResps) != len(requests) {
		return nil, fmt.Errorf("response count mismatch: got %d responses for %d requests",
			len(wireResps), len(requests))
	}

	responses := make([]dtobtc.ImportDescriptorsResponse, len(wireResps))
	for i, r := range wireResps {
		resp := dtobtc.ImportDescriptorsResponse{
			Success:  r.Success,
			Warnings: r.Warnings,
		}
		if r.Error != nil {
			resp.Error = &dtobtc.ImportDescriptorsError{
				Code:    r.Error.Code,
				Message: r.Error.Message,
			}
		}
		responses[i] = resp
	}
	return responses, nil
}

// ImportMulti imports addresses/scripts with optional redeem scripts (legacy wallets).
//
// This is used for BCH and older Bitcoin Core versions that don't support descriptors.
// Allows importing P2SH addresses with redeem scripts for multisig support.
//
// Parameters:
//   - requests: List of import requests
//   - options: Import options (rescan, etc.)
//
// Returns:
//   - List of responses (one per request, in the same order)
//   - Error if the RPC call fails
//
// Notes:
//   - Use "now" timestamp to skip rescanning (fastest)
//   - Set watchonly=true for watch-only wallets
//   - Include redeemscript for P2SH multisig addresses
//
// Reference: Bitcoin Core RPC documentation - importmulti
func (b *Bitcoin) ImportMulti(
	requests []dtobtc.ImportMultiRequest,
	options *dtobtc.ImportMultiOptions,
) ([]dtobtc.ImportMultiResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no addresses to import")
	}

	wireReqs := make([]btcrpc.ImportMultiRequest, len(requests))
	for i, r := range requests {
		wireReqs[i] = btcrpc.ImportMultiRequest{
			ScriptPubKey: r.ScriptPubKey,
			Timestamp:    r.Timestamp,
			RedeemScript: r.RedeemScript,
			PubKeys:      r.PubKeys,
			Keys:         r.Keys,
			Internal:     r.Internal,
			WatchOnly:    r.WatchOnly,
			Label:        r.Label,
		}
	}

	var wireOpts *btcrpc.ImportMultiOptions
	if options != nil {
		wireOpts = &btcrpc.ImportMultiOptions{Rescan: options.Rescan}
	}

	wireResps, err := b.pkgrpc.ImportMulti(wireReqs, wireOpts)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.ImportMulti(): %w", err)
	}

	if len(wireResps) != len(requests) {
		return nil, fmt.Errorf("response count mismatch: got %d responses for %d requests",
			len(wireResps), len(requests))
	}

	responses := make([]dtobtc.ImportMultiResponse, len(wireResps))
	for i, r := range wireResps {
		resp := dtobtc.ImportMultiResponse{
			Success:  r.Success,
			Warnings: r.Warnings,
		}
		if r.Error != nil {
			resp.Error = &dtobtc.ImportMultiError{
				Code:    r.Error.Code,
				Message: r.Error.Message,
			}
		}
		responses[i] = resp
	}
	return responses, nil
}
