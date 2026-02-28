package xrp

import (
	"fmt"
	"strconv"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/protogen"
	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
)

// ToInfraInstructions converts DTO Instructions to infrastructure Instructions.
func ToInfraInstructions(dto *dtoxrp.Instructions) *protogen.Instructions {
	if dto == nil {
		return nil
	}
	return protogen.Instructions_builder{
		Fee:                    dto.Fee,
		MaxFee:                 dto.MaxFee,
		MaxLedgerVersion:       dto.MaxLedgerVersion,
		MaxLedgerVersionOffset: dto.MaxLedgerVersionOffset,
		Sequence:               dto.Sequence,
		SignersCount:           dto.SignersCount,
	}.Build()
}

// ToDTOInstructions converts infrastructure Instructions to DTO Instructions.
func ToDTOInstructions(infra *protogen.Instructions) *dtoxrp.Instructions {
	if infra == nil {
		return nil
	}
	return &dtoxrp.Instructions{
		Fee:                    infra.GetFee(),
		MaxFee:                 infra.GetMaxFee(),
		MaxLedgerVersion:       infra.GetMaxLedgerVersion(),
		MaxLedgerVersionOffset: infra.GetMaxLedgerVersionOffset(),
		Sequence:               infra.GetSequence(),
		SignersCount:           infra.GetSignersCount(),
	}
}

// ToInfraTxInput converts DTO TxInput to infrastructure TxInput.
func ToInfraTxInput(dto *dtoxrp.TxInput) *TxInput {
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
func ToDTOTxInput(infra *TxInput) *dtoxrp.TxInput {
	if infra == nil {
		return nil
	}
	return &dtoxrp.TxInput{
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
func ToDTOSentTx(infra *SentTx) *dtoxrp.SentTx {
	if infra == nil {
		return nil
	}
	return &dtoxrp.SentTx{
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
func ToDTOTxInfo(infra *TxInfo) *dtoxrp.TxInfo {
	if infra == nil {
		return nil
	}

	// Convert balance changes
	balanceChanges := make(map[string][]dtoxrp.TxAmount)
	for key, amounts := range infra.Outcome.BalanceChanges {
		dtoAmounts := make([]dtoxrp.TxAmount, len(amounts))
		for i, amt := range amounts {
			dtoAmounts[i] = dtoxrp.TxAmount{
				Currency: amt.Currency,
				Value:    amt.Value,
			}
		}
		balanceChanges[key] = dtoAmounts
	}

	// Convert orderbook changes
	orderbookChanges := make(map[string][]dtoxrp.TxOrderbookChange)
	for key, changes := range infra.Outcome.OrderbookChanges {
		dtoChanges := make([]dtoxrp.TxOrderbookChange, len(changes))
		for i, change := range changes {
			dtoChanges[i] = dtoxrp.TxOrderbookChange{
				Direction: change.Direction,
				Quantity: dtoxrp.TxAmount{
					Currency: change.Quantity.Currency,
					Value:    change.Quantity.Value,
				},
				TotalPrice: dtoxrp.TxTotalPrice{
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

	return &dtoxrp.TxInfo{
		Type:     infra.Type,
		Address:  infra.Address,
		Sequence: infra.Sequence,
		ID:       infra.ID,
		Specification: dtoxrp.TxSpecification{
			Source: dtoxrp.TxSpecSource{
				Address: infra.Specification.Source.Address,
				MaxAmount: dtoxrp.TxAmount{
					Currency: infra.Specification.Source.MaxAmount.Currency,
					Value:    infra.Specification.Source.MaxAmount.Value,
				},
			},
			Destination: dtoxrp.TxSpecDestination{
				Address: infra.Specification.Destination.Address,
			},
		},
		Outcome: dtoxrp.TxOutcome{
			Result:           infra.Outcome.Result,
			Timestamp:        infra.Outcome.Timestamp.Format("2006-01-02T15:04:05Z"),
			Fee:              infra.Outcome.Fee,
			BalanceChanges:   balanceChanges,
			OrderbookChanges: orderbookChanges,
			LedgerVersion:    infra.Outcome.LedgerVersion,
			IndexInLedger:    infra.Outcome.IndexInLedger,
			DeliveredAmount: dtoxrp.TxAmount{
				Currency: infra.Outcome.DeliveredAmount.Currency,
				Value:    infra.Outcome.DeliveredAmount.Value,
			},
		},
	}
}

// ToInfraXRPKeyType converts DTO XRPKeyType to infrastructure XRPKeyType.
func ToInfraXRPKeyType(dto dtoxrp.XRPKeyType) XRPKeyType {
	return XRPKeyType(dto)
}

// ToDTOXRPKeyType converts infrastructure XRPKeyType to DTO XRPKeyType.
func ToDTOXRPKeyType(infra XRPKeyType) dtoxrp.XRPKeyType {
	return dtoxrp.XRPKeyType(infra)
}

// ToDTOResponseGetAccountInfo converts infrastructure ResponseGetAccountInfo to DTO.
func ToDTOResponseGetAccountInfo(infra *protogen.ResponseGetAccountInfo) *dtoxrp.ResponseGetAccountInfo {
	if infra == nil {
		return nil
	}
	return &dtoxrp.ResponseGetAccountInfo{
		Sequence:                       infra.GetSequence(),
		XrpBalance:                     infra.GetXrpBalance(),
		OwnerCount:                     infra.GetOwnerCount(),
		PreviousAffectingTransactionID: infra.GetPreviousAffectingTransactionID(),
		PreviousAffectingTransactionLedgerVersion: infra.GetPreviousAffectingTransactionLedgerVersion(),
	}
}

// ToDTOResponseGenerateAddress converts infrastructure ResponseGenerateAddress to DTO.
func ToDTOResponseGenerateAddress(infra *protogen.ResponseGenerateAddress) *dtoxrp.ResponseGenerateAddress {
	if infra == nil {
		return nil
	}
	return &dtoxrp.ResponseGenerateAddress{
		XAddress:       infra.GetXAddress(),
		ClassicAddress: infra.GetClassicAddress(),
		Address:        infra.GetAddress(),
		Secret:         infra.GetSecret(),
	}
}

// ToDTOResponseGenerateXAddress converts infrastructure ResponseGenerateXAddress to DTO.
func ToDTOResponseGenerateXAddress(infra *protogen.ResponseGenerateXAddress) *dtoxrp.ResponseGenerateXAddress {
	if infra == nil {
		return nil
	}
	return &dtoxrp.ResponseGenerateXAddress{
		XAddress: infra.GetXAddress(),
		Secret:   infra.GetSecret(),
	}
}

// ToDTOResponseAccountChannelsFromPkg converts pkg ResponseAccountChannels to DTO.
func ToDTOResponseAccountChannelsFromPkg(pkg *xrprpc.ResponseAccountChannels) *dtoxrp.ResponseAccountChannels {
	if pkg == nil {
		return nil
	}

	channels := make([]dtoxrp.AccountChannel, len(pkg.Result.Channels))
	for i, ch := range pkg.Result.Channels {
		channels[i] = dtoxrp.AccountChannel{
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

	return &dtoxrp.ResponseAccountChannels{
		Account:   pkg.Result.Account,
		Channels:  channels,
		Validated: pkg.Result.Validated,
	}
}

// ToDTOResponseAccountInfoFromPkg converts pkg ResponseAccountInfo to DTO.
func ToDTOResponseAccountInfoFromPkg(pkg *xrprpc.ResponseAccountInfo) *dtoxrp.ResponseAccountInfo {
	if pkg == nil {
		return nil
	}
	return &dtoxrp.ResponseAccountInfo{
		Account:            pkg.Result.AccountData.Account,
		Balance:            pkg.Result.AccountData.Balance,
		Sequence:           uint64(pkg.Result.AccountData.Sequence),
		OwnerCount:         uint64(pkg.Result.AccountData.OwnerCount),
		Flags:              uint64(pkg.Result.AccountData.Flags),
		LedgerCurrentIndex: uint64(pkg.Result.LedgerCurrentIndex),
		Validated:          pkg.Result.Validated,
	}
}

// ToDTOResponseServerInfoFromPkg converts pkg ResponseServerInfo to DTO.
func ToDTOResponseServerInfoFromPkg(pkg *xrprpc.ResponseServerInfo) *dtoxrp.ResponseServerInfo {
	if pkg == nil {
		return nil
	}
	return &dtoxrp.ResponseServerInfo{
		BuildVersion:    pkg.Result.Info.BuildVersion,
		CompleteLedgers: pkg.Result.Info.CompleteLedgers,
		HostID:          pkg.Result.Info.Hostid,
		IOLatencyMS:     uint64(pkg.Result.Info.IoLatencyMs),
		LastClose: dtoxrp.ServerLastClose{
			ConvergeTimeS: pkg.Result.Info.LastClose.ConvergeTimeS,
			Proposers:     uint64(pkg.Result.Info.LastClose.Proposers),
		},
		LoadFactor:  uint64(pkg.Result.Info.LoadFactor),
		Peers:       uint64(pkg.Result.Info.Peers),
		PubkeyNode:  pkg.Result.Info.PubkeyNode,
		ServerState: pkg.Result.Info.ServerState,
		ValidatedLedger: dtoxrp.ValidatedLedger{
			Age:            uint64(pkg.Result.Info.ValidatedLedger.Age),
			BaseFeeXRP:     fmt.Sprintf("%f", pkg.Result.Info.ValidatedLedger.BaseFeeXrp),
			Hash:           pkg.Result.Info.ValidatedLedger.Hash,
			ReserveBaseXRP: strconv.Itoa(pkg.Result.Info.ValidatedLedger.ReserveBaseXrp),
			ReserveIncXRP:  strconv.Itoa(pkg.Result.Info.ValidatedLedger.ReserveIncXrp),
			Seq:            uint64(pkg.Result.Info.ValidatedLedger.Seq),
		},
		ValidationQuorum: uint64(pkg.Result.Info.ValidationQuorum),
	}
}

// ToDTOResponseValidationCreateFromPkg converts pkg ResponseValidationCreate to DTO.
func ToDTOResponseValidationCreateFromPkg(pkg *xrprpc.ResponseValidationCreate) *dtoxrp.ResponseValidationCreate {
	if pkg == nil {
		return nil
	}
	return &dtoxrp.ResponseValidationCreate{
		ValidationPublicKey: pkg.Result.ValidationPublicKey,
		ValidationSeed:      pkg.Result.ValidationSeed,
		ValidationKey:       pkg.Result.ValidationKey,
	}
}

// ToDTOResponseWalletProposeFromPkg converts pkg ResponseWalletPropose to DTO.
func ToDTOResponseWalletProposeFromPkg(pkg *xrprpc.ResponseWalletPropose) *dtoxrp.ResponseWalletPropose {
	if pkg == nil {
		return nil
	}

	warning := ""
	if pkg.Status == StatusCodeError.String() {
		warning = pkg.Error
	}

	return &dtoxrp.ResponseWalletPropose{
		MasterSeed:    pkg.Result.MasterSeed,
		MasterSeedHex: pkg.Result.MasterSeedHex,
		MasterKey:     pkg.Result.MasterKey,
		AccountID:     pkg.Result.AccountID,
		PublicKey:     pkg.Result.PublicKey,
		PublicKeyHex:  pkg.Result.PublicKeyHex,
		KeyType:       pkg.Result.KeyType,
		Warning:       warning,
	}
}

// ToDTOSignerListSetTxInput converts infrastructure SignerListSetTxInput to DTO.
func ToDTOSignerListSetTxInput(infra *SignerListSetTxInput) *dtoxrp.SignerListSetTxInput {
	if infra == nil {
		return nil
	}

	entries := make([]dtoxrp.SignerListEntry, len(infra.SignerEntries))
	for i, entry := range infra.SignerEntries {
		entries[i] = dtoxrp.SignerListEntry{
			Account:      entry.SignerEntry.Account,
			SignerWeight: entry.SignerEntry.SignerWeight,
		}
	}

	return &dtoxrp.SignerListSetTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		SignerQuorum:       infra.SignerQuorum,
		SignerEntries:      entries,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOTrustSetTxInput converts infrastructure TrustSetTxInput to DTO.
func ToDTOTrustSetTxInput(infra *TrustSetTxInput) *dtoxrp.TrustSetTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.TrustSetTxInput{
		TransactionType: infra.TransactionType,
		Account:         infra.Account,
		LimitAmount: dtoxrp.IssuedCurrencyAmount{
			Currency: infra.LimitAmount.Currency,
			Issuer:   infra.LimitAmount.Issuer,
			Value:    infra.LimitAmount.Value,
		},
		QualityIn:          infra.QualityIn,
		QualityOut:         infra.QualityOut,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOEscrowCreateTxInput converts infrastructure EscrowCreateTxInput to DTO.
func ToDTOEscrowCreateTxInput(infra *EscrowCreateTxInput) *dtoxrp.EscrowCreateTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.EscrowCreateTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Amount:             infra.Amount,
		Destination:        infra.Destination,
		CancelAfter:        infra.CancelAfter,
		FinishAfter:        infra.FinishAfter,
		Condition:          infra.Condition,
		DestinationTag:     infra.DestinationTag,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOEscrowFinishTxInput converts infrastructure EscrowFinishTxInput to DTO.
func ToDTOEscrowFinishTxInput(infra *EscrowFinishTxInput) *dtoxrp.EscrowFinishTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.EscrowFinishTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Owner:              infra.Owner,
		OfferSequence:      infra.OfferSequence,
		Condition:          infra.Condition,
		Fulfillment:        infra.Fulfillment,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOEscrowCancelTxInput converts infrastructure EscrowCancelTxInput to DTO.
func ToDTOEscrowCancelTxInput(infra *EscrowCancelTxInput) *dtoxrp.EscrowCancelTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.EscrowCancelTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Owner:              infra.Owner,
		OfferSequence:      infra.OfferSequence,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOPaymentChannelCreateTxInput converts infrastructure PaymentChannelCreateTxInput to DTO.
func ToDTOPaymentChannelCreateTxInput(infra *PaymentChannelCreateTxInput) *dtoxrp.PaymentChannelCreateTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.PaymentChannelCreateTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Amount:             infra.Amount,
		Destination:        infra.Destination,
		SettleDelay:        infra.SettleDelay,
		PublicKey:          infra.PublicKey,
		CancelAfter:        infra.CancelAfter,
		DestinationTag:     infra.DestinationTag,
		SourceTag:          infra.SourceTag,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOPaymentChannelFundTxInput converts infrastructure PaymentChannelFundTxInput to DTO.
func ToDTOPaymentChannelFundTxInput(infra *PaymentChannelFundTxInput) *dtoxrp.PaymentChannelFundTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.PaymentChannelFundTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Channel:            infra.Channel,
		Amount:             infra.Amount,
		Expiration:         infra.Expiration,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTOPaymentChannelClaimTxInput converts infrastructure PaymentChannelClaimTxInput to DTO.
func ToDTOPaymentChannelClaimTxInput(infra *PaymentChannelClaimTxInput) *dtoxrp.PaymentChannelClaimTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.PaymentChannelClaimTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		Channel:            infra.Channel,
		Balance:            infra.Balance,
		Amount:             infra.Amount,
		Signature:          infra.Signature,
		PublicKey:          infra.PublicKey,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTONFTokenMintTxInput converts infrastructure NFTokenMintTxInput to DTO.
func ToDTONFTokenMintTxInput(infra *NFTokenMintTxInput) *dtoxrp.NFTokenMintTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.NFTokenMintTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		NFTokenTaxon:       infra.NFTokenTaxon,
		Issuer:             infra.Issuer,
		TransferFee:        infra.TransferFee,
		URI:                infra.URI,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTONFTokenBurnTxInput converts infrastructure NFTokenBurnTxInput to DTO.
func ToDTONFTokenBurnTxInput(infra *NFTokenBurnTxInput) *dtoxrp.NFTokenBurnTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.NFTokenBurnTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		NFTokenID:          infra.NFTokenID,
		Owner:              infra.Owner,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTONFTokenCreateOfferTxInput converts infrastructure NFTokenCreateOfferTxInput to DTO.
func ToDTONFTokenCreateOfferTxInput(infra *NFTokenCreateOfferTxInput) *dtoxrp.NFTokenCreateOfferTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.NFTokenCreateOfferTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		NFTokenID:          infra.NFTokenID,
		Amount:             infra.Amount,
		Owner:              infra.Owner,
		Expiration:         infra.Expiration,
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

// ToDTONFTokenAcceptOfferTxInput converts infrastructure NFTokenAcceptOfferTxInput to DTO.
func ToDTONFTokenAcceptOfferTxInput(infra *NFTokenAcceptOfferTxInput) *dtoxrp.NFTokenAcceptOfferTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.NFTokenAcceptOfferTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		NFTokenSellOffer:   infra.NFTokenSellOffer,
		NFTokenBuyOffer:    infra.NFTokenBuyOffer,
		NFTokenBrokerFee:   infra.NFTokenBrokerFee,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}

// ToDTONFTokenCancelOfferTxInput converts infrastructure NFTokenCancelOfferTxInput to DTO.
func ToDTONFTokenCancelOfferTxInput(infra *NFTokenCancelOfferTxInput) *dtoxrp.NFTokenCancelOfferTxInput {
	if infra == nil {
		return nil
	}

	return &dtoxrp.NFTokenCancelOfferTxInput{
		TransactionType:    infra.TransactionType,
		Account:            infra.Account,
		NFTokenOffers:      infra.NFTokenOffers,
		Fee:                infra.Fee,
		Flags:              infra.Flags,
		LastLedgerSequence: infra.LastLedgerSequence,
		Sequence:           infra.Sequence,
		SigningPubKey:      infra.SigningPubKey,
		TxnSignature:       infra.TxnSignature,
		Hash:               infra.Hash,
	}
}
