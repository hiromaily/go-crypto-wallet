package btc

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
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
func (b *Bitcoin) ImportAddress(pubkey string) error {
	err := b.Client.ImportAddress(pubkey)
	if err != nil {
		return fmt.Errorf("fail to call ImportAddress(): %w", err)
	}

	return nil
}

// ImportAddressWithoutReScan import pubkey without rescan
func (b *Bitcoin) ImportAddressWithoutReScan(pubkey string) error {
	err := b.ImportAddressWithLabel(pubkey, "", false)
	if err != nil {
		return fmt.Errorf("fail to call ImportAddressWithoutReScan(): %w", err)
	}

	return nil
}

// ImportAddressWithLabel import geven address with label to wallet
// - rescan is adjustable
func (b *Bitcoin) ImportAddressWithLabel(address, label string, rescan bool) error {
	bAddress, err := json.Marshal(address)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(address): %w", err)
	}

	// addresses
	bLabel, err := json.Marshal(label)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(label): %w", err)
	}

	// rescan
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return fmt.Errorf("fail to call json.Marchal(rescan): %w", err)
	}
	jsonRawMsg := []json.RawMessage{bAddress, bLabel, bRescan}

	// call importaddress
	_, err = b.Client.RawRequest("importaddress", jsonRawMsg)
	if err != nil {
		return fmt.Errorf("fail to call client.RawRequest(importaddress): %w", err)
	}

	return nil
}
