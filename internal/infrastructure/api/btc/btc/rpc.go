package btc

import (
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// EstimateSmartFee estimates the fee per kilobyte needed for a transaction.
func (b *Bitcoin) EstimateSmartFee(confirmationBlock int) (float64, error) {
	return b.pkgrpc.EstimateSmartFee(confirmationBlock)
}

// GetNetworkInfo returns network information from the connected node.
func (b *Bitcoin) GetNetworkInfo() (*dtobtc.NetworkInfo, error) {
	result, err := b.pkgrpc.GetNetworkInfo()
	if err != nil {
		return nil, fmt.Errorf("fail to call pkgrpc.GetNetworkInfo(): %w", err)
	}
	return toNetworkInfo(result), nil
}

func toNetworkInfo(r *btcrpc.GetNetworkInfoResult) *dtobtc.NetworkInfo {
	networks := make([]dtobtc.NetworkInfoNetwork, len(r.Networks))
	for i, n := range r.Networks {
		networks[i] = dtobtc.NetworkInfoNetwork{
			Name:      n.Name,
			Limited:   n.Limited,
			Reachable: n.Reachable,
			Proxy:     n.Proxy,
		}
	}
	localAddrs := make([]dtobtc.NetworkInfoLocalAddress, len(r.LocalAddresses))
	for i, a := range r.LocalAddresses {
		localAddrs[i] = dtobtc.NetworkInfoLocalAddress{
			Address: a.Address,
			Port:    a.Port,
			Score:   a.Score,
		}
	}
	return &dtobtc.NetworkInfo{
		Version:         r.Version,
		Subversion:      r.Subversion,
		ProtocolVersion: r.ProtocolVersion,
		Connections:     r.Connections,
		Networks:        networks,
		RelayFee:        r.RelayFee,
		LocalAddresses:  localAddrs,
		Warnings:        r.Warnings.Value,
	}
}

// Logging returns logging information from the connected node.
func (b *Bitcoin) Logging() (*dtobtc.LoggingStatus, error) {
	result, err := b.pkgrpc.Logging()
	if err != nil {
		return nil, fmt.Errorf("fail to call pkgrpc.Logging(): %w", err)
	}
	return &dtobtc.LoggingStatus{
		Net:         result.Net,
		Tor:         result.Tor,
		Mempool:     result.Mempool,
		HTTP:        result.HTTP,
		Bench:       result.Bench,
		Zmq:         result.Zmq,
		Walletdb:    result.Walletdb,
		RPC:         result.RPC,
		Estimatefee: result.Estimatefee,
		Addrman:     result.Addrman,
		Selectcoins: result.Selectcoins,
		Reindex:     result.Reindex,
		Cmpctblock:  result.Cmpctblock,
		Rand:        result.Rand,
		Prune:       result.Prune,
		Proxy:       result.Proxy,
		Mempoolrej:  result.Mempoolrej,
		Libevent:    result.Libevent,
		Coindb:      result.Coindb,
		Qt:          result.Qt,
		Leveldb:     result.Leveldb,
		Validation:  result.Validation,
	}, nil
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
