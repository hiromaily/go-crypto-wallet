package eth

import (
	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
)

// ToDomainUserAmount converts infrastructure UserAmount to domain UserAmount.
func ToDomainUserAmount(infra UserAmount) domainEthereum.UserAmount {
	return domainEthereum.UserAmount{
		Address: infra.Address,
		Amount:  infra.Amount,
	}
}

// ToDomainUserAmounts converts a slice of infrastructure UserAmount to domain UserAmount.
func ToDomainUserAmounts(infra []UserAmount) []domainEthereum.UserAmount {
	result := make([]domainEthereum.UserAmount, len(infra))
	for i, ua := range infra {
		result[i] = ToDomainUserAmount(ua)
	}
	return result
}

// FromDomainQuantityTag converts domain QuantityTag to infrastructure QuantityTag.
func FromDomainQuantityTag(domain domainEthereum.QuantityTag) QuantityTag {
	return QuantityTag(domain)
}

// ToDomainQuantityTag converts infrastructure QuantityTag to domain QuantityTag.
func ToDomainQuantityTag(infra QuantityTag) domainEthereum.QuantityTag {
	return domainEthereum.QuantityTag(infra)
}

// ToDomainBlockInfo converts infrastructure BlockInfo to domain BlockInfo.
func ToDomainBlockInfo(infra *BlockInfo) *domainEthereum.BlockInfo {
	if infra == nil {
		return nil
	}
	return &domainEthereum.BlockInfo{
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
func ToDomainResponseGetTransaction(infra *ResponseGetTransaction) *domainEthereum.ResponseGetTransaction {
	if infra == nil {
		return nil
	}
	return &domainEthereum.ResponseGetTransaction{
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
) *domainEthereum.ResponseGetTransactionReceipt {
	if infra == nil {
		return nil
	}
	return &domainEthereum.ResponseGetTransactionReceipt{
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
func ToDomainResponseSyncing(infra *ResponseSyncing) *domainEthereum.ResponseSyncing {
	if infra == nil {
		return nil
	}
	return &domainEthereum.ResponseSyncing{
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
