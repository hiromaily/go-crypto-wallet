package ethtx

import (
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
)

// ToDomainRawTx converts infrastructure RawTx to domain RawTx.
func ToDomainRawTx(infra *RawTx) *domainETH.RawTx {
	if infra == nil {
		return nil
	}
	return &domainETH.RawTx{
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
func FromDomainRawTx(domain *domainETH.RawTx) *RawTx {
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
