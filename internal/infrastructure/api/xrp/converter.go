package xrp

import (
	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
)

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

// ToInfraXRPKeyType converts DTO XRPKeyType to infrastructure XRPKeyType.
func ToInfraXRPKeyType(dto dtoxrp.XRPKeyType) XRPKeyType {
	return XRPKeyType(dto)
}

// ToDTOXRPKeyType converts infrastructure XRPKeyType to DTO XRPKeyType.
func ToDTOXRPKeyType(infra XRPKeyType) dtoxrp.XRPKeyType {
	return dtoxrp.XRPKeyType(infra)
}
