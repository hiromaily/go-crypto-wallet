package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

type signTransactionUseCase struct {
	txSigner       apieth.TxSigner
	txFileRepo     file.TransactionFileRepositorier
	accountKeyRepo repocold.ETHAccountKeyRepositorier
}

// NewSignTransactionUseCase creates a new SignTransactionUseCase for ETH keygen offline signing.
// No network or RPC calls are made during signing (fully offline operation).
func NewSignTransactionUseCase(
	txSigner apieth.TxSigner,
	txFileRepo file.TransactionFileRepositorier,
	accountKeyRepo repocold.ETHAccountKeyRepositorier,
) keygenusecase.SignTransactionUseCase {
	return &signTransactionUseCase{
		txSigner:       txSigner,
		txFileRepo:     txFileRepo,
		accountKeyRepo: accountKeyRepo,
	}
}

func (u *signTransactionUseCase) Sign(
	ctx context.Context,
	input keygenusecase.SignTransactionInput,
) (keygenusecase.SignTransactionOutput, error) {
	// Parse and validate the file path (must be an unsigned transaction file)
	actionType, _, txID, signedCount, err := u.txFileRepo.ValidateFilePath(input.FilePath, domainTx.TxTypeUnsigned)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, err
	}

	// Read unsigned transaction JSON file
	txFile, err := u.txFileRepo.ReadETHJSONFile(input.FilePath)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("failed to read ETH transaction file: %w", err)
	}
	if err = txFile.Validate(); err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("invalid ETH transaction file: %w", err)
	}

	// Retrieve account key for the sender address
	accountKey, err := u.accountKeyRepo.GetByAddress(txFile.From)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf(
			"key not found for address %s: %w", txFile.From, err,
		)
	}
	if accountKey.AccountExtendedPrivkey == "" {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf(
			"accountXpriv not found for address %s: run key generation first", txFile.From,
		)
	}

	// Derive child private key from accountXpriv at the key's BIP-44 index
	// Full derivation: m/44'/60'/account'/0/idx  (accountXpriv is at account level)
	privKey, err := deriveChildPrivKey(accountKey.AccountExtendedPrivkey, uint32(accountKey.Idx))
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf(
			"failed to derive private key at index %d: %w", accountKey.Idx, err,
		)
	}

	// Build RawTx from the transaction file for the signing function
	rawTx := &domainETH.RawTx{
		From:  txFile.From,
		To:    txFile.To,
		TxHex: txFile.RawTxHex,
	}

	// Sign offline using private key — no network calls
	chainID := new(big.Int).SetUint64(txFile.ChainID)
	signedRawTx, err := u.txSigner.SignTxWithPrivateKey(rawTx, privKey, chainID)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Build signed transaction file; leave the original unsigned file unchanged
	signedTxFile := *txFile
	signedTxFile.TxType = string(domainTx.TxTypeSigned)
	signedTxFile.SignedTxHex = signedRawTx.TxHex

	// Write the signed JSON file
	outPath := u.txFileRepo.CreateFilePath(actionType, domainTx.TxTypeSigned, txID, signedCount+1)
	generatedFileName, err := u.txFileRepo.WriteETHJSONFile(outPath, &signedTxFile)
	if err != nil {
		return keygenusecase.SignTransactionOutput{}, fmt.Errorf("failed to write signed transaction file: %w", err)
	}

	return keygenusecase.SignTransactionOutput{
		FilePath:      generatedFileName,
		IsDone:        true,
		SignedCount:   1,
		UnsignedCount: 0,
	}, nil
}

// deriveChildPrivKey derives the child private key from an account-level extended private key
// at the given BIP-44 child index. The account-level key is at m/44'/60'/account'; this
// function derives the external chain (0) and then the child at idx, yielding m/44'/60'/account'/0/idx.
func deriveChildPrivKey(accountXpriv string, idx uint32) (*ecdsa.PrivateKey, error) {
	extKey, err := hdkeychain.NewKeyFromString(accountXpriv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse accountXpriv: %w", err)
	}

	// Derive external chain key (index 0): m/44'/60'/account'/0
	externalKey, err := extKey.Derive(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive external chain key: %w", err)
	}

	// Derive child key at idx: m/44'/60'/account'/0/idx
	childKey, err := externalKey.Derive(idx)
	if err != nil {
		return nil, fmt.Errorf("failed to derive child key at index %d: %w", idx, err)
	}

	ecPrivKey, err := childKey.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get EC private key: %w", err)
	}

	return ecPrivKey.ToECDSA(), nil
}
