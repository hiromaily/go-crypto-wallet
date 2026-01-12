package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ToAddressInfo converts infrastructure GetAddressInfoResult to application AddressInfo
func ToAddressInfo(result *GetAddressInfoResult) *dtobtc.AddressInfo {
	if result == nil {
		return nil
	}

	dto := &dtobtc.AddressInfo{
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
func ToValidateAddressResult(result *ValidateAddressResult) *dtobtc.ValidateAddressResult {
	if result == nil {
		return nil
	}

	return &dtobtc.ValidateAddressResult{
		IsValid:      result.IsValid,
		Address:      result.Address,
		ScriptPubKey: result.ScriptPubKey,
		IsScript:     result.IsScript,
		IsWitness:    result.IsWitness,
	}
}

// ToNetworkInfo converts infrastructure to application DTO
func ToNetworkInfo(result *GetNetworkInfoResult) *dtobtc.NetworkInfo {
	if result == nil {
		return nil
	}

	networks := make([]dtobtc.NetworkAddress, len(result.Networks))
	for i, net := range result.Networks {
		networks[i] = dtobtc.NetworkAddress{
			Name:                      net.Name,
			Limited:                   net.Limited,
			Reachable:                 net.Reachable,
			Proxy:                     net.Proxy,
			ProxyRandomizeCredentials: net.ProxyRandomizeCredentials,
		}
	}

	localAddrs := make([]dtobtc.LocalAddress, len(result.Localaddresses))
	for i, addr := range result.Localaddresses {
		localAddrs[i] = dtobtc.LocalAddress{
			Address: addr.Address,
			Port:    uint16(addr.Port),
			Score:   int32(addr.Score),
		}
	}

	return &dtobtc.NetworkInfo{
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
		Warnings:        result.Warnings.String(),
	}
}

// ToBlockchainInfo converts infrastructure to application DTO
func ToBlockchainInfo(result *GetBlockchainInfoResult) *dtobtc.BlockchainInfo {
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

	return &dtobtc.BlockchainInfo{
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
		Warnings:             result.Warnings.String(),
	}
}

// ToTransactionResult converts infrastructure GetTransactionResult to application DTO
func ToTransactionResult(result *GetTransactionResult, btc *Bitcoin) (*dtobtc.TransactionResult, error) {
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

	details := make([]dtobtc.TransactionDetail, len(result.Details))
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

		details[i] = dtobtc.TransactionDetail{
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

	return &dtobtc.TransactionResult{
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
func ToRawTransaction(result *TxRawResult, btc *Bitcoin) (*dtobtc.RawTransaction, error) {
	if result == nil {
		return nil, nil
	}

	vin := make([]dtobtc.RawTransactionInput, len(result.Vin))
	for i, input := range result.Vin {
		vin[i] = dtobtc.RawTransactionInput{
			TxID: input.Txid,
			Vout: input.Vout,
			ScriptSig: dtobtc.ScriptSig{
				Asm: input.ScriptSig.Asm,
				Hex: input.ScriptSig.Hex,
			},
			Sequence: input.Sequence,
			Witness:  nil, // TxRawVin doesn't have Witness field
		}
	}

	vout := make([]dtobtc.RawTransactionOutput, len(result.Vout))
	for i, output := range result.Vout {
		value, err := btc.FloatToAmount(output.Value)
		if err != nil {
			return nil, err
		}

		address := ""
		if len(output.ScriptPubKey.Addresses) > 0 {
			address = output.ScriptPubKey.Addresses[0]
		}

		vout[i] = dtobtc.RawTransactionOutput{
			Value: value,
			Index: output.N,
			ScriptPubKey: dtobtc.ScriptPubKey{
				Asm:     output.ScriptPubKey.Asm,
				Hex:     output.ScriptPubKey.Hex,
				ReqSigs: int32(output.ScriptPubKey.ReqSigs),
				Type:    output.ScriptPubKey.Type,
				Address: address,
			},
		}
	}

	return &dtobtc.RawTransaction{
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
) *dtobtc.FundRawTransactionResult {
	if result == nil {
		return nil
	}

	// Fee is already in satoshis, convert directly to btcutil.Amount
	fee := btcutil.Amount(result.Fee)

	return &dtobtc.FundRawTransactionResult{
		Hex:       result.Hex,
		Fee:       fee,
		ChangePos: int32(result.Changepos),
	}
}

// FromPreviousTx converts application PreviousTx to infrastructure PrevTx
func FromPreviousTx(prevTxs []dtobtc.PreviousTx, btc *Bitcoin) ([]PrevTx, error) {
	if prevTxs == nil {
		return nil, nil
	}

	result := make([]PrevTx, len(prevTxs))
	for i, tx := range prevTxs {
		amount := float64(tx.Amount) / 1e8 // Convert satoshis to BTC

		result[i] = PrevTx{
			Txid:          tx.TxID,
			Vout:          tx.Vout,
			ScriptPubKey:  tx.ScriptPubKey,
			RedeemScript:  tx.RedeemScript,
			WitnessScript: tx.WitnessScript,
			Amount:        amount,
		}
	}
	return result, nil
}

// ToPreviousTxList converts infrastructure PrevTx slice to application PreviousTx slice
func ToPreviousTxList(prevTxs []PrevTx, btc *Bitcoin) ([]dtobtc.PreviousTx, error) {
	if prevTxs == nil {
		return nil, nil
	}

	result := make([]dtobtc.PreviousTx, len(prevTxs))
	for i, tx := range prevTxs {
		amount, err := btc.FloatToAmount(tx.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to convert amount for prevTx %d: %w", i, err)
		}

		result[i] = dtobtc.PreviousTx{
			TxID:          tx.Txid,
			Vout:          tx.Vout,
			ScriptPubKey:  tx.ScriptPubKey,
			RedeemScript:  tx.RedeemScript,
			WitnessScript: "", // Not available in PrevTx
			Amount:        amount,
		}
	}
	return result, nil
}

// ToUnspentOutput converts infrastructure ListUnspentResult to application DTO
func ToUnspentOutput(result *ListUnspentResult, btc *Bitcoin) (*dtobtc.UnspentOutput, error) {
	if result == nil {
		return nil, nil
	}

	// Convert float64 amount to btcutil.Amount
	amount, err := btc.FloatToAmount(result.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to convert amount: %w", err)
	}

	return &dtobtc.UnspentOutput{
		TxID:          result.TxID,
		Vout:          result.Vout,
		Address:       result.Address,
		Account:       "", // Not provided by Bitcoin Core RPC
		ScriptPubKey:  result.ScriptPubKey,
		Amount:        amount,
		Confirmations: result.Confirmations,
		RedeemScript:  result.RedeemScript,
		WitnessScript: result.WitnessScript, // Bitcoin Core provides this for P2WSH addresses
		Spendable:     result.Spendable,
		Solvable:      result.Solvable,
		Safe:          result.Safe,
		Label:         result.Label,
	}, nil
}

// ToUnspentOutputList converts slice of infrastructure results to DTOs
func ToUnspentOutputList(results []ListUnspentResult, btc *Bitcoin) ([]dtobtc.UnspentOutput, error) {
	if results == nil {
		return nil, nil
	}

	outputs := make([]dtobtc.UnspentOutput, 0, len(results))
	for _, result := range results {
		dto, err := ToUnspentOutput(&result, btc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert unspent output: %w", err)
		}
		if dto != nil {
			outputs = append(outputs, *dto)
		}
	}
	return outputs, nil
}

// FromUnspentOutput converts application UnspentOutput to infrastructure type
func FromUnspentOutput(output *dtobtc.UnspentOutput) *ListUnspentResult {
	if output == nil {
		return nil
	}

	// Convert btcutil.Amount to float64 BTC
	amount := output.Amount.ToBTC()

	return &ListUnspentResult{
		TxID:          output.TxID,
		Vout:          output.Vout,
		Address:       output.Address,
		Label:         output.Label,
		RedeemScript:  output.RedeemScript,
		WitnessScript: output.WitnessScript,
		ScriptPubKey:  output.ScriptPubKey,
		Amount:        amount,
		Confirmations: output.Confirmations,
		Spendable:     output.Spendable,
		Solvable:      output.Solvable,
		Safe:          output.Safe,
	}
}

// ToParsedPSBT converts infrastructure ParsedPSBT to application DTO
//
//nolint:gocyclo // Complex mapping function with many fields
func ToParsedPSBT(infraPSBT *ParsedPSBT, btc *Bitcoin) (*dtobtc.ParsedPSBT, error) {
	if infraPSBT == nil || infraPSBT.Packet == nil {
		return nil, nil
	}

	packet := infraPSBT.Packet
	unsignedTx := packet.UnsignedTx

	// Map transaction
	tx := dtobtc.ParsedPSBTTx{
		TxID:     unsignedTx.TxHash().String(),
		Hash:     unsignedTx.WitnessHash().String(),
		Version:  unsignedTx.Version,
		LockTime: unsignedTx.LockTime,
		Vin:      make([]dtobtc.ParsedPSBTVin, len(unsignedTx.TxIn)),
		Vout:     make([]dtobtc.ParsedPSBTVout, len(unsignedTx.TxOut)),
	}

	// Map inputs (from unsigned tx)
	for i, txIn := range unsignedTx.TxIn {
		tx.Vin[i] = dtobtc.ParsedPSBTVin{
			TxID:     txIn.PreviousOutPoint.Hash.String(),
			Vout:     txIn.PreviousOutPoint.Index,
			Sequence: txIn.Sequence,
		}
	}

	// Map outputs (from unsigned tx)
	for i, txOut := range unsignedTx.TxOut {
		amount := btcutil.Amount(txOut.Value)
		tx.Vout[i] = dtobtc.ParsedPSBTVout{
			Value:        amount,
			ScriptPubKey: hex.EncodeToString(txOut.PkScript),
		}
	}

	// Map PSBT inputs (metadata)
	inputs := make([]dtobtc.ParsedPSBTInput, len(packet.Inputs))
	for i, input := range packet.Inputs {
		parsedInput := dtobtc.ParsedPSBTInput{
			PartialSignatures: make(map[string]string),
			SigHashType:       uint32(input.SighashType),
			RedeemScript:      hex.EncodeToString(input.RedeemScript),
			WitnessScript:     hex.EncodeToString(input.WitnessScript),
			FinalScriptSig:    hex.EncodeToString(input.FinalScriptSig),
			Unknown:           make(map[string]string),
			BIP32Derivation:   make([]dtobtc.BIP32Derivation, 0),
		}

		// Map partial signatures
		for _, partialSig := range input.PartialSigs {
			pubkeyHex := hex.EncodeToString(partialSig.PubKey)
			sigHex := hex.EncodeToString(partialSig.Signature)
			parsedInput.PartialSignatures[pubkeyHex] = sigHex
		}

		// Map final witness - deserialize witness stack
		witnessStack, err := parseWitnessStack(input.FinalScriptWitness)
		if err != nil {
			return nil, fmt.Errorf("failed to parse witness stack for input %d: %w", i, err)
		}
		parsedInput.FinalScriptWitness = witnessStack

		// Map witness UTXO
		if input.WitnessUtxo != nil {
			amount := btcutil.Amount(input.WitnessUtxo.Value)
			parsedInput.WitnessUTXO = &dtobtc.ParsedPSBTUTXO{
				Amount:       amount,
				ScriptPubKey: hex.EncodeToString(input.WitnessUtxo.PkScript),
			}
		}

		// Map non-witness UTXO
		if input.NonWitnessUtxo != nil {
			rawTx := &TxRawResult{
				Txid:     input.NonWitnessUtxo.TxHash().String(),
				Hash:     input.NonWitnessUtxo.WitnessHash().String(),
				Size:     int32(input.NonWitnessUtxo.SerializeSize()),
				Vsize:    int32(input.NonWitnessUtxo.SerializeSize()),
				Weight:   int32(input.NonWitnessUtxo.SerializeSize() * 4),
				Version:  uint32(input.NonWitnessUtxo.Version),
				Locktime: input.NonWitnessUtxo.LockTime,
				Vin:      make([]TxRawVin, len(input.NonWitnessUtxo.TxIn)),
				Vout:     make([]TxRawVout, len(input.NonWitnessUtxo.TxOut)),
			}

			// Map vin
			for j, txIn := range input.NonWitnessUtxo.TxIn {
				rawTx.Vin[j] = TxRawVin{
					Txid: txIn.PreviousOutPoint.Hash.String(),
					Vout: txIn.PreviousOutPoint.Index,
					ScriptSig: ScriptSig{
						Hex: hex.EncodeToString(txIn.SignatureScript),
					},
					Sequence: txIn.Sequence,
				}
			}

			// Map vout
			for j, txOut := range input.NonWitnessUtxo.TxOut {
				rawTx.Vout[j] = TxRawVout{
					Value: float64(txOut.Value) / 1e8,
					N:     uint32(j),
					ScriptPubKey: ScriptPubKey{
						Hex: hex.EncodeToString(txOut.PkScript),
					},
				}
			}

			nonWitnessTx, err := ToRawTransaction(rawTx, btc)
			if err != nil {
				return nil, fmt.Errorf("failed to convert non-witness UTXO: %w", err)
			}
			parsedInput.NonWitnessUTXO = nonWitnessTx
		}

		// Map BIP32 derivation
		for _, deriv := range input.Bip32Derivation {
			// Convert MasterKeyFingerprint (uint32) to 4-byte hex string
			fingerprint := fmt.Sprintf("%08x", deriv.MasterKeyFingerprint)
			parsedInput.BIP32Derivation = append(parsedInput.BIP32Derivation, dtobtc.BIP32Derivation{
				PubKey:      hex.EncodeToString(deriv.PubKey),
				MasterKeyID: fingerprint,
				Path:        derivationPathToString(deriv.Bip32Path),
			})
		}

		// Map unknown fields
		for _, unknown := range input.Unknowns {
			parsedInput.Unknown[hex.EncodeToString(unknown.Key)] = hex.EncodeToString(unknown.Value)
		}

		inputs[i] = parsedInput
	}

	// Map PSBT outputs (metadata)
	outputs := make([]dtobtc.ParsedPSBTOutput, len(packet.Outputs))
	for i, output := range packet.Outputs {
		parsedOutput := dtobtc.ParsedPSBTOutput{
			RedeemScript:    hex.EncodeToString(output.RedeemScript),
			WitnessScript:   hex.EncodeToString(output.WitnessScript),
			Unknown:         make(map[string]string),
			BIP32Derivation: make([]dtobtc.BIP32Derivation, 0),
		}

		// Map BIP32 derivation
		for _, deriv := range output.Bip32Derivation {
			// Convert MasterKeyFingerprint (uint32) to 4-byte hex string
			fingerprint := fmt.Sprintf("%08x", deriv.MasterKeyFingerprint)
			parsedOutput.BIP32Derivation = append(parsedOutput.BIP32Derivation, dtobtc.BIP32Derivation{
				PubKey:      hex.EncodeToString(deriv.PubKey),
				MasterKeyID: fingerprint,
				Path:        derivationPathToString(deriv.Bip32Path),
			})
		}

		// Map unknown fields
		for _, unknown := range output.Unknowns {
			parsedOutput.Unknown[hex.EncodeToString(unknown.Key)] = hex.EncodeToString(unknown.Value)
		}

		outputs[i] = parsedOutput
	}

	// Calculate fee: total input amount - total output amount
	var totalInput btcutil.Amount
	for idx, input := range packet.Inputs {
		if input.WitnessUtxo != nil {
			// SegWit transactions use WitnessUtxo
			totalInput += btcutil.Amount(input.WitnessUtxo.Value)
		} else if input.NonWitnessUtxo != nil && idx < len(unsignedTx.TxIn) {
			// Legacy transactions use NonWitnessUtxo
			// Get the amount from the corresponding output in the previous transaction
			vout := unsignedTx.TxIn[idx].PreviousOutPoint.Index
			if int(vout) < len(input.NonWitnessUtxo.TxOut) {
				totalInput += btcutil.Amount(input.NonWitnessUtxo.TxOut[vout].Value)
			}
		}
	}

	var totalOutput btcutil.Amount
	for _, txOut := range unsignedTx.TxOut {
		totalOutput += btcutil.Amount(txOut.Value)
	}

	fee := totalInput - totalOutput

	// Map unknown global fields
	unknownGlobal := make(map[string]string)
	for _, unknown := range packet.Unknowns {
		unknownGlobal[hex.EncodeToString(unknown.Key)] = hex.EncodeToString(unknown.Value)
	}

	return &dtobtc.ParsedPSBT{
		Tx:         tx,
		Unknown:    unknownGlobal,
		Inputs:     inputs,
		Outputs:    outputs,
		Fee:        fee,
		IsComplete: infraPSBT.IsComplete,
	}, nil
}

// derivationPathToString converts BIP32 path to string format
func derivationPathToString(path []uint32) string {
	if len(path) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("m")
	for _, component := range path {
		hardened := component >= 0x80000000
		if hardened {
			sb.WriteString(fmt.Sprintf("/%d'", component-0x80000000))
		} else {
			sb.WriteString(fmt.Sprintf("/%d", component))
		}
	}
	return sb.String()
}

// parseWitnessStack deserializes a witness stack from its byte representation.
// The witness stack format is:
//   - VarInt: number of witness items
//   - For each item:
//   - VarInt: length of the witness item
//   - []byte: the witness item data
func parseWitnessStack(witnessData []byte) ([]string, error) {
	if len(witnessData) == 0 {
		return []string{}, nil
	}

	reader := bytes.NewReader(witnessData)

	// Read number of witness items
	count, err := wire.ReadVarInt(reader, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to read witness count: %w", err)
	}

	result := make([]string, 0, count)
	for i := uint64(0); i < count; i++ {
		// Read length of this witness item
		length, err := wire.ReadVarInt(reader, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to read witness item %d length: %w", i, err)
		}

		// Read the witness item data
		item := make([]byte, length)
		n, err := reader.Read(item)
		if err != nil {
			return nil, fmt.Errorf("failed to read witness item %d: %w", i, err)
		}
		if uint64(n) != length {
			return nil, fmt.Errorf("witness item %d: expected %d bytes, read %d", i, length, n)
		}

		result = append(result, hex.EncodeToString(item))
	}

	return result, nil
}

// ToLoggingResult converts infrastructure LoggingResult to application DTO
func ToLoggingResult(result *LoggingResult) *dtobtc.LoggingResult {
	if result == nil {
		return nil
	}

	return &dtobtc.LoggingResult{
		Net:         result.Net,
		TorControl:  result.Tor,
		MemPool:     result.Mempool,
		RPC:         result.RPC,
		EstimateFee: result.Estimatefee,
		Bench:       result.Bench,
		// Map Zmq to ZMQPubRawBlock and ZMQPubRawTx (both set to same value)
		ZMQPubRawBlock: result.Zmq,
		ZMQPubRawTx:    result.Zmq,
		Coindb:         result.Coindb,
		SelectCoins:    result.Selectcoins,
	}
}

// ToMultisigAddress converts infrastructure AddMultisigAddressResult to application DTO
func ToMultisigAddress(result *AddMultisigAddressResult) *dtobtc.MultisigAddress {
	if result == nil {
		return nil
	}

	return &dtobtc.MultisigAddress{
		Address:      result.Address,
		RedeemScript: result.RedeemScript,
		Descriptor:   "", // Not provided by infrastructure result
	}
}

// FromAddressType converts Bitcoin Core RPC AddressType to user-facing AddrType
// Note: Converts "bech32m" (Bitcoin Core value) to "taproot" (user-facing value)
func FromAddressType(addrType domainBitcoin.AddressType) domainAddress.AddrType {
	switch addrType {
	case domainBitcoin.AddressTypeLegacy:
		return domainAddress.AddrTypeLegacy
	case domainBitcoin.AddressTypeP2SHSegwit:
		return domainAddress.AddrTypeP2shSegwit
	case domainBitcoin.AddressTypeBech32:
		return domainAddress.AddrTypeBech32
	case domainBitcoin.AddressTypeTaproot: // "bech32m" -> "taproot"
		return domainAddress.AddrTypeTaproot
	default:
		// This indicates a programming error (unhandled domain type).
		// Panicking helps to catch this issue during development.
		panic(fmt.Sprintf("unhandled domain address type: %s", addrType))
	}
}

// ToAddressType converts user-facing AddrType to Bitcoin Core RPC AddressType
// Note: Converts "taproot" (user-facing value) to "bech32m" (Bitcoin Core value)
func ToAddressType(addrType domainAddress.AddrType) domainBitcoin.AddressType {
	switch addrType {
	case domainAddress.AddrTypeLegacy:
		return domainBitcoin.AddressTypeLegacy
	case domainAddress.AddrTypeP2shSegwit:
		return domainBitcoin.AddressTypeP2SHSegwit
	case domainAddress.AddrTypeBech32:
		return domainBitcoin.AddressTypeBech32
	case domainAddress.AddrTypeTaproot: // "taproot" -> "bech32m"
		return domainBitcoin.AddressTypeTaproot
	case domainAddress.AddrTypeBCHCashAddr:
		// BCH uses legacy address format in domain model
		return domainBitcoin.AddressTypeLegacy
	case domainAddress.AddrTypeETH:
		// ETH not supported in Bitcoin domain
		return domainBitcoin.AddressTypeLegacy
	default:
		// Log warning for unknown infrastructure types to improve observability
		logger.Warn("unknown infrastructure address type, defaulting to legacy",
			"address_type", addrType)
		return domainBitcoin.AddressTypeLegacy
	}
}

// ToBTCVersion converts infrastructure BTCVersion to domain Version
func ToBTCVersion(ver BTCVersion) domainBitcoin.Version {
	return domainBitcoin.Version(ver)
}

// FromBTCVersion converts domain Version to infrastructure BTCVersion
func FromBTCVersion(ver domainBitcoin.Version) BTCVersion {
	return BTCVersion(ver)
}
