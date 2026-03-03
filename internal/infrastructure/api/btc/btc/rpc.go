package btc

import (
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// EstimateSmartFee estimates the fee per kilobyte needed for a transaction.
func (b *Bitcoin) EstimateSmartFee(confirmationBlock int) (float64, error) {
	return b.pkgrpc.EstimateSmartFee(confirmationBlock)
}

// GetNetworkInfo returns network information from the connected node.
func (b *Bitcoin) GetNetworkInfo() (*btcrpc.GetNetworkInfoResult, error) {
	return b.pkgrpc.GetNetworkInfo()
}

// Logging returns logging information from the connected node.
func (b *Bitcoin) Logging() (*btcrpc.LoggingResult, error) {
	return b.pkgrpc.Logging()
}

// EncryptWallet encrypts the wallet with the given passphrase.
func (b *Bitcoin) EncryptWallet(passphrase string) error {
	return b.pkgrpc.EncryptWallet(passphrase)
}

// WalletLock removes the wallet encryption key from memory, locking the wallet.
func (b *Bitcoin) WalletLock() error {
	return b.pkgrpc.WalletLock()
}

// WalletPassphrase stores the wallet decryption key in memory for the given timeout.
func (b *Bitcoin) WalletPassphrase(passphrase string, timeoutSecs int64) error {
	return b.pkgrpc.WalletPassphrase(passphrase, timeoutSecs)
}

// WalletPassphraseChange changes the wallet passphrase.
func (b *Bitcoin) WalletPassphraseChange(old, newPass string) error {
	return b.pkgrpc.WalletPassphraseChange(old, newPass)
}

// DumpWallet dumps all wallet keys to a server-side file.
func (b *Bitcoin) DumpWallet(fileName string) error {
	return b.pkgrpc.DumpWallet(fileName)
}

// ImportWallet imports keys from a wallet dump file.
func (b *Bitcoin) ImportWallet(fileName string) error {
	return b.pkgrpc.ImportWallet(fileName)
}

// GetAddressesByLabelMap returns the raw RPC label-to-purpose map without address decoding.
// Bitcoin.GetAddressesByLabel decodes results into btcutil.Address using BTC chain config,
// which BCH cannot use. BCH's GetAddressesByLabel override calls this to get the raw
// address strings and apply BCH-specific decoding.
func (b *Bitcoin) GetAddressesByLabelMap(labelName string) (map[string]btcrpc.Purpose, error) {
	return b.pkgrpc.GetAddressesByLabel(labelName)
}
