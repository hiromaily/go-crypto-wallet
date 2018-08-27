package api

import (
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
