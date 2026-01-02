package xrp

import (
	"fmt"
	"strconv"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
)

// ToInfraInstructions converts DTO Instructions to infrastructure Instructions.
func ToInfraInstructions(dto *dtoRipple.Instructions) *Instructions {
	if dto == nil {
		return nil
	}
	return &Instructions{
		Fee:                    dto.Fee,
		MaxFee:                 dto.MaxFee,
		MaxLedgerVersion:       dto.MaxLedgerVersion,
		MaxLedgerVersionOffset: dto.MaxLedgerVersionOffset,
		Sequence:               dto.Sequence,
		SignersCount:           dto.SignersCount,
	}
}

// ToDTOInstructions converts infrastructure Instructions to DTO Instructions.
func ToDTOInstructions(infra *Instructions) *dtoRipple.Instructions {
	if infra == nil {
		return nil
	}
	return &dtoRipple.Instructions{
		Fee:                    infra.Fee,
		MaxFee:                 infra.MaxFee,
		MaxLedgerVersion:       infra.MaxLedgerVersion,
		MaxLedgerVersionOffset: infra.MaxLedgerVersionOffset,
		Sequence:               infra.Sequence,
		SignersCount:           infra.SignersCount,
	}
}

// ToInfraTxInput converts DTO TxInput to infrastructure TxInput.
func ToInfraTxInput(dto *dtoRipple.TxInput) *TxInput {
	if dto == nil {
		return nil
	}
	return &TxInput{
		TransactionType:    dto.TransactionType,
		Account:            dto.Account,
		Amount:             dto.Amount,
		Destination:        dto.Destination,
		Fee:                dto.Fee,
		Flags:              dto.Flags,
		LastLedgerSequence: dto.LastLedgerSequence,
		Sequence:           dto.Sequence,
		SigningPubKey:      dto.SigningPubKey,
		TxnSignature:       dto.TxnSignature,
		Hash:               dto.Hash,
	}
}

// ToDTOTxInput converts infrastructure TxInput to DTO TxInput.
func ToDTOTxInput(infra *TxInput) *dtoRipple.TxInput {
	if infra == nil {
		return nil
	}
	return &dtoRipple.TxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Amount:             infra.Amount,
		Destination:        infra.Destination,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOSentTx converts infrastructure SentTx to DTO SentTx.
func ToDTOSentTx(infra *SentTx) *dtoRipple.SentTx {
	if infra == nil {
		return nil
	}
	return &dtoRipple.SentTx{
		ResultCode:          infra.ResultCode,
		ResultMessage:       infra.ResultMessage,
		EngineResult:        infra.EngineResult,
		EngineResultCode:    infra.EngineResultCode,
		EngineResultMessage: infra.EngineResultMessage,
		TxBlob:              infra.TxBlob,
		TxJSON:              *ToDTOTxInput(&infra.TxJSON),
	}
}

// ToDTOTxInfo converts infrastructure TxInfo to DTO TxInfo.
func ToDTOTxInfo(infra *TxInfo) *dtoRipple.TxInfo {
	if infra == nil {
		return nil
	}

	// Convert balance changes
	balanceChanges := make(map[string][]dtoRipple.TxAmount)
	for key, amounts := range infra.Outcome.BalanceChanges {
		dtoAmounts := make([]dtoRipple.TxAmount, len(amounts))
		for i, amt := range amounts {
			dtoAmounts[i] = dtoRipple.TxAmount{
				Currency: amt.Currency,
				Value:    amt.Value,
			}
		}
		balanceChanges[key] = dtoAmounts
	}

	// Convert orderbook changes
	orderbookChanges := make(map[string][]dtoRipple.TxOrderbookChange)
	for key, changes := range infra.Outcome.OrderbookChanges {
		dtoChanges := make([]dtoRipple.TxOrderbookChange, len(changes))
		for i, change := range changes {
			dtoChanges[i] = dtoRipple.TxOrderbookChange{
				Direction: change.Direction,
				Quantity: dtoRipple.TxAmount{
					Currency: change.Quantity.Currency,
					Value:    change.Quantity.Value,
				},
				TotalPrice: dtoRipple.TxTotalPrice{
					Currency:     change.TotalPrice.Currency,
					Counterparty: change.TotalPrice.Counterparty,
					Value:        change.TotalPrice.Value,
				},
				Sequence:          change.Sequence,
				Status:            change.Status,
				MakerExchangeRate: change.MakerExchangeRate,
			}
		}
		orderbookChanges[key] = dtoChanges
	}

	return &dtoRipple.TxInfo{
		Type:     infra.Type,
		Address:  infra.Address,
		Sequence: infra.Sequence,
		ID:       infra.ID,
		Specification: dtoRipple.TxSpecification{
			Source: dtoRipple.TxSpecSource{
				Address: infra.Specification.Source.Address,
				MaxAmount: dtoRipple.TxAmount{
					Currency: infra.Specification.Source.MaxAmount.Currency,
					Value:    infra.Specification.Source.MaxAmount.Value,
				},
			},
			Destination: dtoRipple.TxSpecDestination{
				Address: infra.Specification.Destination.Address,
			},
		},
		Outcome: dtoRipple.TxOutcome{
			Result:           infra.Outcome.Result,
			Timestamp:        infra.Outcome.Timestamp.Format("2006-01-02T15:04:05Z"),
			Fee:              infra.Outcome.Fee,
			BalanceChanges:   balanceChanges,
			OrderbookChanges: orderbookChanges,
			LedgerVersion:    infra.Outcome.LedgerVersion,
			IndexInLedger:    infra.Outcome.IndexInLedger,
			DeliveredAmount: dtoRipple.TxAmount{
				Currency: infra.Outcome.DeliveredAmount.Currency,
				Value:    infra.Outcome.DeliveredAmount.Value,
			},
		},
	}
}

// ToInfraXRPKeyType converts DTO XRPKeyType to infrastructure XRPKeyType.
func ToInfraXRPKeyType(dto dtoRipple.XRPKeyType) XRPKeyType {
	return XRPKeyType(dto)
}

// ToDTOXRPKeyType converts infrastructure XRPKeyType to DTO XRPKeyType.
func ToDTOXRPKeyType(infra XRPKeyType) dtoRipple.XRPKeyType {
	return dtoRipple.XRPKeyType(infra)
}

// ToDTOResponseGetAccountInfo converts infrastructure ResponseGetAccountInfo to DTO.
func ToDTOResponseGetAccountInfo(infra *ResponseGetAccountInfo) *dtoRipple.ResponseGetAccountInfo {
	if infra == nil {
		return nil
	}
	return &dtoRipple.ResponseGetAccountInfo{
		Sequence:                       infra.Sequence,
		XrpBalance:                     infra.XrpBalance,
		OwnerCount:                     infra.OwnerCount,
		PreviousAffectingTransactionID: infra.PreviousAffectingTransactionID,
		PreviousAffectingTransactionLedgerVersion: infra.PreviousAffectingTransactionLedgerVersion,
	}
}

// ToDTOResponseGenerateAddress converts infrastructure ResponseGenerateAddress to DTO.
func ToDTOResponseGenerateAddress(infra *ResponseGenerateAddress) *dtoRipple.ResponseGenerateAddress {
	if infra == nil {
		return nil
	}
	return &dtoRipple.ResponseGenerateAddress{
		XAddress:       infra.XAddress,
		ClassicAddress: infra.ClassicAddress,
		Address:        infra.Address,
		Secret:         infra.Secret,
	}
}

// ToDTOResponseGenerateXAddress converts infrastructure ResponseGenerateXAddress to DTO.
func ToDTOResponseGenerateXAddress(infra *ResponseGenerateXAddress) *dtoRipple.ResponseGenerateXAddress {
	if infra == nil {
		return nil
	}
	return &dtoRipple.ResponseGenerateXAddress{
		XAddress: infra.XAddress,
		Secret:   infra.Secret,
	}
}

// ToDTOResponseAccountChannels converts infrastructure ResponseAccountChannels to DTO.
func ToDTOResponseAccountChannels(infra *ResponseAccountChannels) *dtoRipple.ResponseAccountChannels {
	if infra == nil {
		return nil
	}

	channels := make([]dtoRipple.AccountChannel, len(infra.Result.Channels))
	for i, ch := range infra.Result.Channels {
		channels[i] = dtoRipple.AccountChannel{
			ChannelID:      ch.ChannelID,
			Account:        ch.Account,
			Destination:    ch.DestinationAccount,
			Amount:         ch.Amount,
			Balance:        ch.Balance,
			SettleDelay:    uint64(ch.SettleDelay),
			PublicKey:      ch.PublicKey,
			DestinationTag: uint32(ch.DestinationTag),
			CancelAfter:    uint64(ch.CancelAfter),
			Expiration:     uint64(ch.Expiration),
		}
	}

	return &dtoRipple.ResponseAccountChannels{
		Account:   infra.Result.Account,
		Channels:  channels,
		Validated: infra.Result.Validated,
	}
}

// ToDTOResponseAccountInfo converts infrastructure ResponseAccountInfo to DTO.
func ToDTOResponseAccountInfo(infra *ResponseAccountInfo) *dtoRipple.ResponseAccountInfo {
	if infra == nil {
		return nil
	}
	return &dtoRipple.ResponseAccountInfo{
		Account:            infra.Result.AccountData.Account,
		Balance:            infra.Result.AccountData.Balance,
		Sequence:           uint64(infra.Result.AccountData.Sequence),
		OwnerCount:         uint64(infra.Result.AccountData.OwnerCount),
		Flags:              uint64(infra.Result.AccountData.Flags),
		LedgerCurrentIndex: uint64(infra.Result.LedgerCurrentIndex),
		Validated:          infra.Result.Validated,
	}
}

// ToDTOResponseServerInfo converts infrastructure ResponseServerInfo to DTO.
func ToDTOResponseServerInfo(infra *ResponseServerInfo) *dtoRipple.ResponseServerInfo {
	if infra == nil {
		return nil
	}
	return &dtoRipple.ResponseServerInfo{
		BuildVersion:    infra.Result.Info.BuildVersion,
		CompleteLedgers: infra.Result.Info.CompleteLedgers,
		HostID:          infra.Result.Info.Hostid,
		IOLatencyMS:     uint64(infra.Result.Info.IoLatencyMs),
		LastClose: dtoRipple.ServerLastClose{
			ConvergeTimeS: infra.Result.Info.LastClose.ConvergeTimeS,
			Proposers:     uint64(infra.Result.Info.LastClose.Proposers),
		},
		LoadFactor:  uint64(infra.Result.Info.LoadFactor),
		Peers:       uint64(infra.Result.Info.Peers),
		PubkeyNode:  infra.Result.Info.PubkeyNode,
		ServerState: infra.Result.Info.ServerState,
		ValidatedLedger: dtoRipple.ValidatedLedger{
			Age:            uint64(infra.Result.Info.ValidatedLedger.Age),
			BaseFeeXRP:     fmt.Sprintf("%f", infra.Result.Info.ValidatedLedger.BaseFeeXrp),
			Hash:           infra.Result.Info.ValidatedLedger.Hash,
			ReserveBaseXRP: strconv.Itoa(infra.Result.Info.ValidatedLedger.ReserveBaseXrp),
			ReserveIncXRP:  strconv.Itoa(infra.Result.Info.ValidatedLedger.ReserveIncXrp),
			Seq:            uint64(infra.Result.Info.ValidatedLedger.Seq),
		},
		ValidationQuorum: uint64(infra.Result.Info.ValidationQuorum),
	}
}

// ToDTOResponseValidationCreate converts infrastructure ResponseValidationCreate to DTO.
func ToDTOResponseValidationCreate(infra *ResponseValidationCreate) *dtoRipple.ResponseValidationCreate {
	if infra == nil {
		return nil
	}
	return &dtoRipple.ResponseValidationCreate{
		ValidationPublicKey: infra.Result.ValidationPublicKey,
		ValidationSeed:      infra.Result.ValidationSeed,
		ValidationKey:       infra.Result.ValidationKey,
	}
}

// ToDTOResponseWalletPropose converts infrastructure ResponseWalletPropose to DTO.
func ToDTOResponseWalletPropose(infra *ResponseWalletPropose) *dtoRipple.ResponseWalletPropose {
	if infra == nil {
		return nil
	}

	warning := ""
	if infra.Status == StatusCodeError.String() {
		warning = infra.Error
	}

	return &dtoRipple.ResponseWalletPropose{
		MasterSeed:    infra.Result.MasterSeed,
		MasterSeedHex: infra.Result.MasterSeedHex,
		MasterKey:     infra.Result.MasterKey,
		AccountID:     infra.Result.AccountID,
		PublicKey:     infra.Result.PublicKey,
		PublicKeyHex:  infra.Result.PublicKeyHex,
		KeyType:       infra.Result.KeyType,
		Warning:       warning,
	}
}
