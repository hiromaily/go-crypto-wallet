package api

import (
	"encoding/json"

	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// ImportPrivKey WIFをwalletに登録する
func (b *Bitcoin) ImportPrivKey(privKeyWIF *btcutil.WIF) error {
	err := b.client.ImportPrivKey(privKeyWIF)
	if err != nil {
		return errors.Errorf("ImportPrivKey(): error: %v", err)
	}

	return nil
}

// ImportPrivKeyLabel WIFをwalletに登録する、label付き
func (b *Bitcoin) ImportPrivKeyLabel(privKeyWIF *btcutil.WIF, label string) error {
	err := b.client.ImportPrivKeyLabel(privKeyWIF, label)
	if err != nil {
		return errors.Errorf("ImportPrivKeyLabel(): error: %v", err)
	}

	return nil
}

// ImportPrivKeyWithoutReScan WIFをwalletに登録する、rescanはしない
func (b *Bitcoin) ImportPrivKeyWithoutReScan(privKeyWIF *btcutil.WIF, label string) error {
	err := b.client.ImportPrivKeyRescan(privKeyWIF, label, false)
	if err != nil {
		return errors.Errorf("ImportPrivKeyRescan(): error: %v", err)
	}

	return nil
}

// ImportAddress アドレスをwalletにimportする
func (b *Bitcoin) ImportAddress(pubkey string) error {
	err := b.client.ImportAddress(pubkey)
	if err != nil {
		return errors.Errorf("ImportAddress(): error: %v", err)
	}

	return nil
}

// ImportAddressWithoutReScan Rescanせずにアドレスをwalletにimportする
func (b *Bitcoin) ImportAddressWithoutReScan(pubkey string) error {
	//FIXME: これはI/Fがずれてて使えない
	//err := b.client.ImportAddressRescan(pubkey, false)
	err := b.ImportAddressWithLabel(pubkey, "", false)
	if err != nil {
		return errors.Errorf("ImportAddressWithoutReScan(): error: %v", err)
	}

	return nil
}

// ImportAddressWithLabel address, label, rescanをパラメータとしaddressをwalletにimportする
func (b *Bitcoin) ImportAddressWithLabel(address, label string, rescan bool) error {
	//requiredSigs
	bAddress, err := json.Marshal(address)
	if err != nil {
		return errors.Errorf("json.Marchal(address): error: %v", err)
	}

	//addresses
	bLabel, err := json.Marshal(label)
	if err != nil {
		return errors.Errorf("json.Marchal(label): error: %v", err)
	}

	//rescan
	bRescan, err := json.Marshal(rescan)
	if err != nil {
		return errors.Errorf("json.Marchal(rescan): error: %v", err)
	}
	jsonRawMsg := []json.RawMessage{bAddress, bLabel, bRescan}

	//call importaddress
	_, err = b.client.RawRequest("importaddress", jsonRawMsg)
	if err != nil {
		return errors.Errorf("client.RawRequest(importaddress): error: %v", err)
	}

	return nil
}
