package eth

import (
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// ToDomainUserAmount converts infrastructure UserAmount to domain UserAmount.
func ToDomainUserAmount(infra UserAmount) domainETH.UserAmount {
	return domainETH.UserAmount{
		Address: infra.Address,
		Amount:  infra.Amount,
	}
}

// ToDomainUserAmounts converts a slice of infrastructure UserAmount to domain UserAmount.
func ToDomainUserAmounts(infra []UserAmount) []domainETH.UserAmount {
	result := make([]domainETH.UserAmount, len(infra))
	for i, ua := range infra {
		result[i] = ToDomainUserAmount(ua)
	}
	return result
}

// ToDomainTransactionReceiptFromPkg converts pkg TransactionReceipt to domain.
func ToDomainTransactionReceiptFromPkg(pkg *ethrpc.TransactionReceipt) *domainETH.TransactionReceipt {
	if pkg == nil {
		return nil
	}
	return &domainETH.TransactionReceipt{
		TransactionHash:   pkg.TransactionHash,
		TransactionIndex:  pkg.TransactionIndex,
		BlockHash:         pkg.BlockHash,
		BlockNumber:       pkg.BlockNumber,
		From:              pkg.From,
		To:                pkg.To,
		CumulativeGasUsed: pkg.CumulativeGasUsed,
		GasUsed:           pkg.GasUsed,
		ContractAddress:   pkg.ContractAddress,
		Status:            pkg.Status,
	}
}
