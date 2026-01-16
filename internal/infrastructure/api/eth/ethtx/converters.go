package ethtx

import (
	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
)

// ToDomainRawTx converts infrastructure RawTx to domain RawTx.
func ToDomainRawTx(infra *RawTx) *domainEthereum.RawTx {
	if infra == nil {
		return nil
	}
	return &domainEthereum.RawTx{
		UUID:  infra.UUID,
		From:  infra.From,
		To:    infra.To,
		Value: infra.Value,
		Nonce: infra.Nonce,
		TxHex: infra.TxHex,
		Hash:  infra.Hash,
	}
}

// FromDomainRawTx converts domain RawTx to infrastructure RawTx.
func FromDomainRawTx(domain *domainEthereum.RawTx) *RawTx {
	if domain == nil {
		return nil
	}
	return &RawTx{
		UUID:  domain.UUID,
		From:  domain.From,
		To:    domain.To,
		Value: domain.Value,
		Nonce: domain.Nonce,
		TxHex: domain.TxHex,
		Hash:  domain.Hash,
	}
}
