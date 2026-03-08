package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// Compile-time check: TransactionFileRepository must satisfy MultisigFileRepositorier.
var _ portfile.MultisigFileRepositorier = (*TransactionFileRepository)(nil)

// CreateMultisigFilePath returns the canonical file path for an ETH multisig transaction file.
// The path has the form: {r.filePath}{actionType}_multisig_{uuid}_{signedCount}.json
func (r *TransactionFileRepository) CreateMultisigFilePath(
	actionType domainTx.ActionType,
	uuid string,
	signedCount int,
) string {
	return fmt.Sprintf("%s%s_multisig_%s_%d.json", r.filePath, actionType.String(), uuid, signedCount)
}

// WriteETHMultisigJSONFile serialises f to indented JSON and writes it to path.
// Unlike WriteETHJSONFile, the path is complete (including .json extension) because
// the file name is deterministically derived from the UUID and signedCount; no timestamp suffix is added.
// Returns the written path on success.
func (r *TransactionFileRepository) WriteETHMultisigJSONFile(
	path string,
	f *dtoeth.ETHMultisigTransactionFile,
) (string, error) {
	if f == nil {
		return "", errors.New("ETHMultisigTransactionFile must not be nil")
	}

	if err := f.Validate(); err != nil {
		return "", fmt.Errorf("invalid ETH multisig transaction file: %w", err)
	}

	cleanPath := filepath.Clean(path)
	if err := r.checkPathTraversal(cleanPath); err != nil {
		return "", err
	}

	if !strings.HasSuffix(strings.ToLower(cleanPath), ".json") {
		return "", fmt.Errorf("multisig file path must have .json extension: %s", cleanPath)
	}

	if err := createDirForFile(cleanPath); err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal ETH multisig transaction: %w", err)
	}

	// 0o600: transaction files contain sensitive multisig data (addresses, hashes, signatures)
	if err = os.WriteFile(cleanPath, jsonData, 0o600); err != nil {
		return "", fmt.Errorf("failed to write ETH multisig transaction file %s: %w", cleanPath, err)
	}

	return cleanPath, nil
}

// ReadETHMultisigJSONFile reads the JSON file at path, deserialises it into
// ETHMultisigTransactionFile, validates it, and returns the result.
func (r *TransactionFileRepository) ReadETHMultisigJSONFile(
	path string,
) (*dtoeth.ETHMultisigTransactionFile, error) {
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		return nil, fmt.Errorf("invalid multisig file extension: %s (expected .json)", path)
	}

	cleanPath := filepath.Clean(path)
	if err := r.checkPathTraversal(cleanPath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ETH multisig transaction file %s: %w", cleanPath, err)
	}

	var f dtoeth.ETHMultisigTransactionFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("failed to parse ETH multisig transaction file %s: %w", cleanPath, err)
	}

	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ETH multisig transaction file data in %s: %w", cleanPath, err)
	}

	return &f, nil
}
