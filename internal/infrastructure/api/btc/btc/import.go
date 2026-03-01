package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// ImportPrivKey import privKey to wallet
// - Rescan  *bool `jsonrpcdefault:"true"`
func (b *Bitcoin) ImportPrivKey(privKeyWIF *btcutil.WIF) error {
	err := b.Client.ImportPrivKey(privKeyWIF)
	if err != nil {
		return fmt.Errorf("fail to call client.ImportPrivKey(): %w", err)
	}

	return nil
}

// ImportPrivKeyLabel import privKey with label to wallet
// - Rescan  *bool `jsonrpcdefault:"true"`
func (b *Bitcoin) ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error {
	err := b.Client.ImportPrivKeyLabel(privKeyWIF, label)
	if err != nil {
		return fmt.Errorf("fail to call client.ImportPrivKeyLabel(): %w", err)
	}

	return nil
}

// ImportPrivKeyWithoutReScan import privKey without rescan to wallet
func (b *Bitcoin) ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error {
	err := b.Client.ImportPrivKeyRescan(privKeyWIF, label, false)
	if err != nil {
		return fmt.Errorf("fail to call ImportPrivKeyRescan(): %w", err)
	}

	return nil
}

// ImportAddress import pubkey to wallet
// Note: This is a legacy wallet method. For descriptor wallets (Bitcoin Core v23.0+),
// consider using importdescriptors instead for better functionality and future compatibility.
func (b *Bitcoin) ImportAddress(pubkey string) error {
	err := b.Client.ImportAddress(pubkey)
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
	if err := b.RPC.ImportAddress(address, label, rescan); err != nil {
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
	requests []btcrpc.ImportDescriptorsRequest,
) ([]btcrpc.ImportDescriptorsResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no descriptors to import")
	}

	responses, err := b.RPC.ImportDescriptors(requests)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.ImportDescriptors(): %w", err)
	}

	if len(responses) != len(requests) {
		return nil, fmt.Errorf("response count mismatch: got %d responses for %d requests",
			len(responses), len(requests))
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
	requests []btcrpc.ImportMultiRequest,
	options *btcrpc.ImportMultiOptions,
) ([]btcrpc.ImportMultiResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no addresses to import")
	}

	responses, err := b.RPC.ImportMulti(requests, options)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.ImportMulti(): %w", err)
	}

	if len(responses) != len(requests) {
		return nil, fmt.Errorf("response count mismatch: got %d responses for %d requests",
			len(responses), len(requests))
	}

	return responses, nil
}
