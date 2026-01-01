package btc

import (
	bitcoindto "github.com/hiromaily/go-crypto-wallet/internal/application/dto/bitcoin"
)

// ToAddressInfo converts infrastructure GetAddressInfoResult to application AddressInfo
func ToAddressInfo(result *GetAddressInfoResult) *bitcoindto.AddressInfo {
	if result == nil {
		return nil
	}

	dto := &bitcoindto.AddressInfo{
		Address:      result.Address,
		ScriptPubKey: result.ScriptPubKey,
		IsWitness:    result.Iswitness,
		Labels:       result.Labels,
		IsCompressed: result.Iscompressed,
		IsChange:     result.Ischange,
		Timestamp:    result.Timestamp,
		IsMine:       result.Ismine,
		IsWatchOnly:  result.Iswatchonly,
		IsScript:     result.Isscript,
		PubKey:       result.Pubkey,
	}

	return dto
}

// ToValidateAddressResult converts infrastructure to application DTO
func ToValidateAddressResult(result *ValidateAddressResult) *bitcoindto.ValidateAddressResult {
	if result == nil {
		return nil
	}

	return &bitcoindto.ValidateAddressResult{
		IsValid:      result.IsValid,
		Address:      result.Address,
		ScriptPubKey: result.ScriptPubKey,
		IsScript:     result.IsScript,
		IsWitness:    result.IsWitness,
	}
}

// ToNetworkInfo converts infrastructure to application DTO
func ToNetworkInfo(result *GetNetworkInfoResult) *bitcoindto.NetworkInfo {
	if result == nil {
		return nil
	}

	networks := make([]bitcoindto.NetworkAddress, len(result.Networks))
	for i, net := range result.Networks {
		networks[i] = bitcoindto.NetworkAddress{
			Name:                      net.Name,
			Limited:                   net.Limited,
			Reachable:                 net.Reachable,
			Proxy:                     net.Proxy,
			ProxyRandomizeCredentials: net.ProxyRandomizeCredentials,
		}
	}

	localAddrs := make([]bitcoindto.LocalAddress, len(result.Localaddresses))
	for i, addr := range result.Localaddresses {
		localAddrs[i] = bitcoindto.LocalAddress{
			Address: addr.Address,
			Port:    uint16(addr.Port),
			Score:   int32(addr.Score),
		}
	}

	return &bitcoindto.NetworkInfo{
		Version:         int32(result.Version),
		SubVersion:      result.Subversion,
		ProtocolVersion: int32(result.Protocolversion),
		LocalServices:   result.Localservices,
		LocalRelay:      result.Localrelay,
		TimeOffset:      result.Timeoffset,
		NetworkActive:   result.Networkactive,
		Connections:     int32(result.Connections),
		Networks:        networks,
		RelayFee:        result.Relayfee,
		IncrementalFee:  result.Incrementalfee,
		LocalAddresses:  localAddrs,
		Warnings:        result.Warnings,
	}
}

// ToBlockchainInfo converts infrastructure to application DTO
func ToBlockchainInfo(result *GetBlockchainInfoResult) *bitcoindto.BlockchainInfo {
	if result == nil {
		return nil
	}

	softForks := make(map[string]any)
	addFork := func(name string, fork Fork) {
		if fork.Type != "" {
			softForks[name] = map[string]any{
				"type":   fork.Type,
				"active": fork.Active,
				"height": fork.Height,
			}
		}
	}

	addFork("bip34", result.SoftForks.Bip34)
	addFork("bip66", result.SoftForks.Bip66)
	addFork("bip65", result.SoftForks.Bip65)
	addFork("csv", result.SoftForks.Csv)
	addFork("segwit", result.SoftForks.Segwit)

	return &bitcoindto.BlockchainInfo{
		Chain:                result.Chain.String(),
		Blocks:               int64(result.Blocks),
		Headers:              int64(result.Headers),
		BestBlockHash:        result.Bestblockhash,
		Difficulty:           result.Difficulty,
		MedianTime:           int64(result.Mediantime),
		VerificationProgress: result.Verificationprogress,
		InitialBlockDownload: result.Initialblockdownload,
		Pruned:               result.Pruned,
		SoftForks:            softForks,
		Warnings:             result.Warnings,
	}
}
