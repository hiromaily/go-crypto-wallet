package eth

import (
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
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

// FromDomainQuantityTag converts domain QuantityTag to infrastructure QuantityTag.
func FromDomainQuantityTag(domain domainETH.QuantityTag) QuantityTag {
	return QuantityTag(domain)
}

// ToDomainQuantityTag converts infrastructure QuantityTag to domain QuantityTag.
func ToDomainQuantityTag(infra QuantityTag) domainETH.QuantityTag {
	return domainETH.QuantityTag(infra)
}

// ToDomainBlockInfo converts infrastructure BlockInfo to domain BlockInfo.
func ToDomainBlockInfo(infra *BlockInfo) *domainETH.BlockInfo {
	if infra == nil {
		return nil
	}
	return &domainETH.BlockInfo{
		Number:           infra.Number,
		Hash:             infra.Hash,
		ParentHash:       infra.ParentHash,
		Nonce:            infra.Nonce,
		Sha3Uncles:       infra.Sha3Uncles,
		LogsBloom:        infra.LogsBloom,
		TransactionsRoot: infra.TransactionsRoot,
		StateRoot:        infra.StateRoot,
		Miner:            infra.Miner,
		Difficulty:       infra.Difficulty,
		TotalDifficulty:  infra.TotalDifficulty,
		ExtraData:        infra.ExtraData,
		Size:             infra.Size,
		GasLimit:         infra.GasLimit,
		GasUsed:          infra.GasUsed,
		Timestamp:        infra.Timestamp,
		Transactions:     infra.Transactions,
		Uncles:           infra.Uncles,
		BaseFeePerGas:    infra.BaseFeePerGas,
	}
}

// ToDomainResponseGetTransaction converts infrastructure ResponseGetTransaction to domain.
func ToDomainResponseGetTransaction(infra *ResponseGetTransaction) *domainETH.ResponseGetTransaction {
	if infra == nil {
		return nil
	}
	return &domainETH.ResponseGetTransaction{
		BlockHash:        infra.BlockHash,
		BlockNumber:      infra.BlockNumber,
		From:             infra.From,
		Gas:              infra.Gas,
		GasPrice:         infra.GasPrice,
		Hash:             infra.Hash,
		Input:            infra.Input,
		Nonce:            infra.Nonce,
		To:               infra.To,
		TransactionIndex: infra.TransactionIndex,
		Value:            infra.Value,
		V:                infra.V,
		R:                infra.R,
		S:                infra.S,
	}
}

// ToDomainResponseGetTransactionReceipt converts infrastructure ResponseGetTransactionReceipt to domain.
func ToDomainResponseGetTransactionReceipt(
	infra *ResponseGetTransactionReceipt,
) *domainETH.ResponseGetTransactionReceipt {
	if infra == nil {
		return nil
	}
	return &domainETH.ResponseGetTransactionReceipt{
		TransactionHash:   infra.TransactionHash,
		TransactionIndex:  infra.TransactionIndex,
		BlockHash:         infra.BlockHash,
		BlockNumber:       infra.BlockNumber,
		From:              infra.From,
		To:                infra.To,
		CumulativeGasUsed: infra.CumulativeGasUsed,
		GasUsed:           infra.GasUsed,
		ContractAddress:   infra.ContractAddress,
		Logs:              infra.Logs,
		LogsBloom:         infra.LogsBloom,
		Status:            infra.Status,
	}
}

// ToDomainResponseSyncing converts infrastructure ResponseSyncing to domain.
func ToDomainResponseSyncing(infra *ResponseSyncing) *domainETH.ResponseSyncing {
	if infra == nil {
		return nil
	}
	return &domainETH.ResponseSyncing{
		StartingBlock:       infra.StartingBlock,
		HighestBlock:        infra.HighestBlock,
		CurrentBlock:        infra.CurrentBlock,
		SyncedAccounts:      infra.SyncedAccounts,
		SyncedAccountBytes:  infra.SyncedAccountBytes,
		SyncedBytecodes:     infra.SyncedBytecodes,
		SyncedBytecodeBytes: infra.SyncedBytecodeBytes,
		SyncedStorage:       infra.SyncedStorage,
		SyncedStorageBytes:  infra.SyncedStorageBytes,
		HealingBytecode:     infra.HealingBytecode,
		HealedBytecodes:     infra.HealedBytecodes,
		HealedBytecodeBytes: infra.HealedBytecodeBytes,
		HealingTrienodes:    infra.HealingTrienodes,
		HealedTrienodes:     infra.HealedTrienodes,
		HealedTrienodeBytes: infra.HealedTrienodeBytes,
	}
}
