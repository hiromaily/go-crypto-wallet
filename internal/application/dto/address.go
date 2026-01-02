package dto

import (
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// AddressFormat is address csv format
type AddressFormat struct {
	CoinTypeCode      domainCoin.CoinTypeCode
	AccountType       domainAccount.AccountType
	P2PKHAddress      string
	P2SHSegwitAddress string
	Bech32Address     string
	TaprootAddress    string
	FullPublicKey     string
	MultisigAddress   string
	Idx               string
}
