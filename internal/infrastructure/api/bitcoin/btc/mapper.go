package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

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

// ToTransactionResult converts infrastructure GetTransactionResult to application DTO
func ToTransactionResult(result *GetTransactionResult, btc *Bitcoin) (*bitcoindto.TransactionResult, error) {
	if result == nil {
		return nil, nil
	}

	amount, err := btc.FloatToAmount(result.Amount)
	if err != nil {
		return nil, err
	}

	fee, err := btc.FloatToAmount(result.Fee)
	if err != nil {
		return nil, err
	}

	details := make([]bitcoindto.TransactionDetail, len(result.Details))
	for i, detail := range result.Details {
		detailAmount, err := btc.FloatToAmount(detail.Amount)
		if err != nil {
			return nil, err
		}

		var detailFee *btcutil.Amount
		if detail.Fee != 0 {
			feeAmt, err := btc.FloatToAmount(detail.Fee)
			if err != nil {
				return nil, err
			}
			detailFee = &feeAmt
		}

		details[i] = bitcoindto.TransactionDetail{
			Address:   detail.Address,
			Category:  detail.Category,
			Amount:    detailAmount,
			Label:     detail.Label,
			Vout:      detail.Vout,
			Fee:       detailFee,
			Abandoned: detail.Abandoned,
		}
	}

	walletConflicts := make([]string, len(result.Walletconflicts))
	for i, conflict := range result.Walletconflicts {
		str, ok := conflict.(string)
		if !ok {
			return nil, fmt.Errorf(
				"unexpected type for wallet conflict at index %d: got %T, want string",
				i, conflict,
			)
		}
		walletConflicts[i] = str
	}

	return &bitcoindto.TransactionResult{
		Amount:          amount,
		Fee:             fee,
		Confirmations:   int64(result.Confirmations),
		BlockHash:       result.Blockhash,
		BlockIndex:      int64(result.Blockindex),
		BlockTime:       int64(result.Blocktime),
		TxID:            result.Txid,
		WalletConflicts: walletConflicts,
		Time:            result.Time,
		TimeReceived:    result.Timereceived,
		Details:         details,
		Hex:             result.Hex,
	}, nil
}

// ToRawTransaction converts infrastructure TxRawResult to application DTO
func ToRawTransaction(result *TxRawResult, btc *Bitcoin) (*bitcoindto.RawTransaction, error) {
	if result == nil {
		return nil, nil
	}

	vin := make([]bitcoindto.RawTransactionInput, len(result.Vin))
	for i, input := range result.Vin {
		vin[i] = bitcoindto.RawTransactionInput{
			TxID: input.Txid,
			Vout: input.Vout,
			ScriptSig: bitcoindto.ScriptSig{
				Asm: input.ScriptSig.Asm,
				Hex: input.ScriptSig.Hex,
			},
			Sequence: input.Sequence,
			Witness:  nil, // TxRawVin doesn't have Witness field
		}
	}

	vout := make([]bitcoindto.RawTransactionOutput, len(result.Vout))
	for i, output := range result.Vout {
		value, err := btc.FloatToAmount(output.Value)
		if err != nil {
			return nil, err
		}

		address := ""
		if len(output.ScriptPubKey.Addresses) > 0 {
			address = output.ScriptPubKey.Addresses[0]
		}

		vout[i] = bitcoindto.RawTransactionOutput{
			Value: value,
			Index: output.N,
			ScriptPubKey: bitcoindto.ScriptPubKey{
				Asm:     output.ScriptPubKey.Asm,
				Hex:     output.ScriptPubKey.Hex,
				ReqSigs: int32(output.ScriptPubKey.ReqSigs),
				Type:    output.ScriptPubKey.Type,
				Address: address,
			},
		}
	}

	return &bitcoindto.RawTransaction{
		Hex:           "", // Need to be set by caller
		TxID:          result.Txid,
		Hash:          result.Hash,
		Size:          result.Size,
		VSize:         result.Vsize,
		Weight:        result.Weight,
		Version:       int32(result.Version),
		LockTime:      result.Locktime,
		Vin:           vin,
		Vout:          vout,
		BlockHash:     "", // TxRawResult doesn't have block info
		Confirmations: 0,  // TxRawResult doesn't have confirmations
		Time:          0,  // TxRawResult doesn't have time
		BlockTime:     0,  // TxRawResult doesn't have blocktime
	}, nil
}

// ToFundRawTransactionResult converts infrastructure to application DTO
func ToFundRawTransactionResult(
	result *FundRawTransactionResult,
) *bitcoindto.FundRawTransactionResult {
	if result == nil {
		return nil
	}

	// Fee is already in satoshis, convert directly to btcutil.Amount
	fee := btcutil.Amount(result.Fee)

	return &bitcoindto.FundRawTransactionResult{
		Hex:       result.Hex,
		Fee:       fee,
		ChangePos: int32(result.Changepos),
	}
}

// FromPreviousTx converts application PreviousTx to infrastructure PrevTx
func FromPreviousTx(prevTxs []bitcoindto.PreviousTx, btc *Bitcoin) ([]PrevTx, error) {
	if prevTxs == nil {
		return nil, nil
	}

	result := make([]PrevTx, len(prevTxs))
	for i, tx := range prevTxs {
		amount := float64(tx.Amount) / 1e8 // Convert satoshis to BTC

		result[i] = PrevTx{
			Txid:         tx.TxID,
			Vout:         tx.Vout,
			ScriptPubKey: tx.ScriptPubKey,
			RedeemScript: tx.RedeemScript,
			Amount:       amount,
		}
	}
	return result, nil
}
