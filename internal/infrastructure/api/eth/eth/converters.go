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

// ToDomainResponseSyncingFromPkg converts pkg ResponseSyncing to domain ResponseSyncing.
func ToDomainResponseSyncingFromPkg(pkg *ethrpc.ResponseSyncing) *domainETH.ResponseSyncing {
	if pkg == nil {
		return nil
	}
	return &domainETH.ResponseSyncing{
		StartingBlock:       pkg.StartingBlock,
		HighestBlock:        pkg.HighestBlock,
		CurrentBlock:        pkg.CurrentBlock,
		SyncedAccounts:      pkg.SyncedAccounts,
		SyncedAccountBytes:  pkg.SyncedAccountBytes,
		SyncedBytecodes:     pkg.SyncedBytecodes,
		SyncedBytecodeBytes: pkg.SyncedBytecodeBytes,
		SyncedStorage:       pkg.SyncedStorage,
		SyncedStorageBytes:  pkg.SyncedStorageBytes,
		HealingBytecode:     pkg.HealingBytecode,
		HealedBytecodes:     pkg.HealedBytecodes,
		HealedBytecodeBytes: pkg.HealedBytecodeBytes,
		HealingTrienodes:    pkg.HealingTrienodes,
		HealedTrienodes:     pkg.HealedTrienodes,
		HealedTrienodeBytes: pkg.HealedTrienodeBytes,
	}
}

// ToDomainBlockInfoFromPkg converts pkg BlockInfo to domain BlockInfo.
func ToDomainBlockInfoFromPkg(pkg *ethrpc.BlockInfo) *domainETH.BlockInfo {
	if pkg == nil {
		return nil
	}
	return &domainETH.BlockInfo{
		Number:           pkg.Number,
		Hash:             pkg.Hash,
		ParentHash:       pkg.ParentHash,
		Nonce:            pkg.Nonce,
		Sha3Uncles:       pkg.Sha3Uncles,
		LogsBloom:        pkg.LogsBloom,
		TransactionsRoot: pkg.TransactionsRoot,
		StateRoot:        pkg.StateRoot,
		Miner:            pkg.Miner,
		Difficulty:       pkg.Difficulty,
		TotalDifficulty:  pkg.TotalDifficulty,
		ExtraData:        pkg.ExtraData,
		Size:             pkg.Size,
		GasLimit:         pkg.GasLimit,
		GasUsed:          pkg.GasUsed,
		Timestamp:        pkg.Timestamp,
		Transactions:     pkg.Transactions,
		Uncles:           pkg.Uncles,
		BaseFeePerGas:    pkg.BaseFeePerGas,
	}
}

// ToDomainResponseGetTransactionFromPkg converts pkg ResponseGetTransaction to domain.
func ToDomainResponseGetTransactionFromPkg(
	pkg *ethrpc.ResponseGetTransaction,
) *domainETH.ResponseGetTransaction {
	if pkg == nil {
		return nil
	}
	return &domainETH.ResponseGetTransaction{
		BlockHash:        pkg.BlockHash,
		BlockNumber:      pkg.BlockNumber,
		From:             pkg.From,
		Gas:              pkg.Gas,
		GasPrice:         pkg.GasPrice,
		Hash:             pkg.Hash,
		Input:            pkg.Input,
		Nonce:            pkg.Nonce,
		To:               pkg.To,
		TransactionIndex: pkg.TransactionIndex,
		Value:            pkg.Value,
		V:                pkg.V,
		R:                pkg.R,
		S:                pkg.S,
	}
}

// ToDomainResponseGetTransactionReceiptFromPkg converts pkg ResponseGetTransactionReceipt to domain.
func ToDomainResponseGetTransactionReceiptFromPkg(
	pkg *ethrpc.ResponseGetTransactionReceipt,
) *domainETH.ResponseGetTransactionReceipt {
	if pkg == nil {
		return nil
	}
	return &domainETH.ResponseGetTransactionReceipt{
		TransactionHash:   pkg.TransactionHash,
		TransactionIndex:  pkg.TransactionIndex,
		BlockHash:         pkg.BlockHash,
		BlockNumber:       pkg.BlockNumber,
		From:              pkg.From,
		To:                pkg.To,
		CumulativeGasUsed: pkg.CumulativeGasUsed,
		GasUsed:           pkg.GasUsed,
		ContractAddress:   pkg.ContractAddress,
		Logs:              pkg.Logs,
		LogsBloom:         pkg.LogsBloom,
		Status:            pkg.Status,
	}
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
