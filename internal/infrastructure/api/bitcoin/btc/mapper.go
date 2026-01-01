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

	// Convert SoftForks struct to map[string]any
	softForks := make(map[string]any)
	if result.SoftForks.Bip34.Type != "" {
		softForks["bip34"] = map[string]any{
			"type":   result.SoftForks.Bip34.Type,
			"active": result.SoftForks.Bip34.Active,
			"height": result.SoftForks.Bip34.Height,
		}
	}
	if result.SoftForks.Bip66.Type != "" {
		softForks["bip66"] = map[string]any{
			"type":   result.SoftForks.Bip66.Type,
			"active": result.SoftForks.Bip66.Active,
			"height": result.SoftForks.Bip66.Height,
		}
	}
	if result.SoftForks.Bip65.Type != "" {
		softForks["bip65"] = map[string]any{
			"type":   result.SoftForks.Bip65.Type,
			"active": result.SoftForks.Bip65.Active,
			"height": result.SoftForks.Bip65.Height,
		}
	}
	if result.SoftForks.Csv.Type != "" {
		softForks["csv"] = map[string]any{
			"type":   result.SoftForks.Csv.Type,
			"active": result.SoftForks.Csv.Active,
			"height": result.SoftForks.Csv.Height,
		}
	}
	if result.SoftForks.Segwit.Type != "" {
		softForks["segwit"] = map[string]any{
			"type":   result.SoftForks.Segwit.Type,
			"active": result.SoftForks.Segwit.Active,
			"height": result.SoftForks.Segwit.Height,
		}
	}

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
