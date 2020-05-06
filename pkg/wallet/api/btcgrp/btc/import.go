package btc

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// ImportPrivKey import privKey to wallet
// - Rescan  *bool `jsonrpcdefault:"true"`
func (b *Bitcoin) ImportPrivKey(privKeyWIF *btcutil.WIF) error {
	err := b.Client.ImportPrivKey(privKeyWIF)
	if err != nil {
		return errors.Wrap(err, "fail to call client.ImportPrivKey()")
	}

	return nil
}

// ImportPrivKeyLabel import privKey with label to wallet
// - Rescan  *bool `jsonrpcdefault:"true"`
func (b *Bitcoin) ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error {
	err := b.Client.ImportPrivKeyLabel(privKeyWIF, label)
	if err != nil {
		return errors.Wrap(err, "fail to call client.ImportPrivKeyLabel()")
	}

	return nil
}

// ImportPrivKeyWithoutReScan import privKey without rescan to wallet
func (b *Bitcoin) ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error {
	err := b.Client.ImportPrivKeyRescan(privKeyWIF, label, false)
	if err != nil {
		return errors.Wrap(err, "fail to call ImportPrivKeyRescan()")
	}

	return nil
}

// ImportAddress import pubkey to wallet
func (b *Bitcoin) ImportAddress(pubkey string) error {
	err := b.Client.ImportAddress(pubkey)
	if err != nil {
		return errors.Wrap(err, "fail to call ImportAddress()")
	}

	return nil
}

// ImportAddressWithoutReScan import pubkey without rescan
func (b *Bitcoin) ImportAddressWithoutReScan(pubkey string) error {
	err := b.ImportAddressWithLabel(pubkey, "", false)
	if err != nil {
		return errors.Wrap(err, "fail to call ImportAddressWithoutReScan()")
	}

	return nil
}

// ImportAddressWithLabel import geven address with label to wallet
// - rescan is adjustable
func (b *Bitcoin) ImportAddressWithLabel(address, label string, rescan bool) error {
	bAddress, err := json.Marshal(address)
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marchal(address)")
	}

	//addresses
	bLabel, err := json.Marshal(label)
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marchal(label)")
	}

	//rescan
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return errors.Wrap(err, "fail to call json.Marchal(rescan)")
	}
	jsonRawMsg := []json.RawMessage{bAddress, bLabel, bRescan}

	//call importaddress
	_, err = b.Client.RawRequest("importaddress", jsonRawMsg)
	if err != nil {
		return errors.Wrap(err, "fail to call client.RawRequest(importaddress)")
	}

	return nil
}
