package eth

import (
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

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
