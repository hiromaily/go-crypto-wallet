package file

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// TransactionFileRepositorier is file storager for tx info
type TransactionFileRepositorier interface {
	CreateFilePath(actionType domainTx.ActionType, txType domainTx.TxType, txID int64, signedCount int) string
	GetFileNameType(filePath string) (*FileName, error)
	ValidateFilePath(
		filePath string,
		expectedTxType domainTx.TxType,
	) (domainTx.ActionType, domainTx.TxType, int64, int, error)
	ReadFile(path string) (string, error)
	ReadFileSlice(path string) ([]string, error)
	WriteFile(path, hexTx string) (string, error)
	WriteFileSlice(path string, data []string) (string, error)

	// PSBT-specific methods (BIP174)
	ReadPSBTFile(path string) (string, error)
	WritePSBTFile(path, psbtBase64 string) (string, error)
}

// TransactionFileRepository is to store transaction info as csv file
type TransactionFileRepository struct {
	filePath string
}

// FileName is object for items in fine name
type FileName struct {
	ActionType  domainTx.ActionType
	TxType      domainTx.TxType
	TxID        int64
	SignedCount int
}

// NewTransactionFileRepository returns TransactionFileRepository
func NewTransactionFileRepository(filePath string) *TransactionFileRepository {
	return &TransactionFileRepository{
		filePath: filePath,
	}
}

// about file structure
// e.g. ./data/tx/deposit/deposit_8_unsigned_0_1534744535097796209.psbt
//  - ./data/tx/ dir : file path
//  - deposit/   dir : actionType
//  - deposit_8_unsigned_0_1534744535097796209.psbt : {actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt

// CreateFilePath create file path for transaction file (PSBT format with .psbt extension)
func (r *TransactionFileRepository) CreateFilePath(
	actionType domainTx.ActionType, txType domainTx.TxType, txID int64, signedCount int,
) string {
	// ./data/tx/deposit/deposit_8_unsigned_0_1534744535097796209.psbt
	// baseDir := fmt.Sprintf("%s%s/", r.filePath, actionType.String())

	// ./data/tx/eth/deposit_8_unsigned_0_1534744535097796209.psbt
	baseDir := r.filePath
	return fmt.Sprintf("%s%s_%d_%s_%d_", baseDir, actionType.String(), txID, txType, signedCount)
}

// GetFileNameType returns as FileName type
func (*TransactionFileRepository) GetFileNameType(filePath string) (*FileName, error) {
	// just file path or full path
	// ./data/tx/deposit/deposit_8_unsigned_0_1534744535097796209.psbt
	tmp := strings.Split(filePath, "/")
	fileName := tmp[len(tmp)-1]

	// Strip .psbt extension if present
	fileName = strings.TrimSuffix(fileName, ".psbt")

	// deposit_5_unsigned_0_1534466246366489473
	// s[0]: actionType
	// s[1]: txID
	// s[2]: txType
	// s[3]: signedCount , first value is 0
	// s[4]: timestamp
	s := strings.Split(fileName, "_")
	if len(s) != 5 {
		return nil, fmt.Errorf("invalid file path: %s", fileName)
	}

	fileNameType := FileName{}

	// Action
	if !domainTx.ValidateActionType(s[0]) {
		return nil, fmt.Errorf("invalid file name: %s", fileName)
	}
	fileNameType.ActionType = domainTx.ActionType(s[0])

	// txID
	var err error
	fileNameType.TxID, err = strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file name: %s", fileName)
	}

	// txType
	if !domainTx.ValidateTxType(s[2]) {
		return nil, fmt.Errorf("invalid name: %s", fileName)
	}
	fileNameType.TxType = domainTx.TxType(s[2])

	// signedCount
	signedCount, err := strconv.Atoi(s[3])
	if err != nil {
		return nil, fmt.Errorf("invalid name: %s", fileName)
	}
	fileNameType.SignedCount = signedCount

	return &fileNameType, nil
}

// ValidateFilePath validate file path which could be full path
func (r *TransactionFileRepository) ValidateFilePath(
	filePath string, expectedTxType domainTx.TxType,
) (domainTx.ActionType, domainTx.TxType, int64, int, error) {
	fileType, err := r.GetFileNameType(filePath)
	if err != nil {
		return "", "", 0, 0, err
	}
	// txType
	// if !(fileType.TxType).Search(expectedTxTypes) {
	if fileType.TxType != expectedTxType {
		return "", "", 0, 0, fmt.Errorf("txType is invalid: %s", fileType.TxType)
	}
	return fileType.ActionType, fileType.TxType, fileType.TxID, fileType.SignedCount, nil
}

// ReadFile read file
func (*TransactionFileRepository) ReadFile(path string) (string, error) {
	ret, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.ReadFile(%s): %w", path, err)
	}

	return string(ret), nil
}

// ReadFileSlice read file for slice
func (*TransactionFileRepository) ReadFileSlice(path string) ([]string, error) {
	file, err := os.Open(path) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("fail to open file: %s: %w", path, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()
	data := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("fail to scan file: %w", err)
	}
	return data, nil
}

// WriteFile write file
func (r *TransactionFileRepository) WriteFile(path, hexTx string) (string, error) {
	// crate directory if not exisiting
	r.createDir(path)

	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName := path + ts

	byteTx := []byte(hexTx)
	err := os.WriteFile(fileName, byteTx, 0o644)
	if err != nil {
		return "", fmt.Errorf("fail to call os.WriteFile(%s): %w", fileName, err)
	}

	return fileName, nil
}

// WriteFileSlice write slice to file
func (r *TransactionFileRepository) WriteFileSlice(path string, data []string) (string, error) {
	// crate directory if not exisiting
	r.createDir(path)

	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName := path + ts

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.OpenFile(%s): %w", fileName, err)
	}
	writer := bufio.NewWriter(file)

	for _, d := range data {
		_, _ = writer.WriteString(d + "\n")
	}

	if err = writer.Flush(); err != nil {
		return "", err
	}

	if err = file.Close(); err != nil {
		return "", fmt.Errorf("failed to close file: %w", err)
	}

	return fileName, nil
}

// ReadPSBTFile reads a base64-encoded PSBT from file
func (*TransactionFileRepository) ReadPSBTFile(path string) (string, error) {
	// Validate extension (case-insensitive)
	if !strings.HasSuffix(strings.ToLower(path), ".psbt") {
		return "", fmt.Errorf("invalid PSBT file extension: %s (expected .psbt)", path)
	}

	// Read file
	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("failed to read PSBT file %s: %w", path, err)
	}

	psbtBase64 := string(data)

	// Validate base64 format
	if _, err := base64.StdEncoding.DecodeString(psbtBase64); err != nil {
		return "", fmt.Errorf("invalid base64 content in PSBT file %s: %w", path, err)
	}

	return psbtBase64, nil
}

// WritePSBTFile writes a base64-encoded PSBT to file with .psbt extension
func (r *TransactionFileRepository) WritePSBTFile(path, psbtBase64 string) (string, error) {
	// Validate base64 format
	if _, err := base64.StdEncoding.DecodeString(psbtBase64); err != nil {
		return "", fmt.Errorf("invalid base64 PSBT data: %w", err)
	}

	// Create directory if not existing
	r.createDir(path)

	// Add timestamp and .psbt extension
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName := path + ts + ".psbt"

	// Write base64 PSBT
	bytePSBT := []byte(psbtBase64)
	err := os.WriteFile(fileName, bytePSBT, 0o644)
	if err != nil {
		return "", fmt.Errorf("failed to write PSBT file %s: %w", fileName, err)
	}

	return fileName, nil
}

func (*TransactionFileRepository) createDir(path string) {
	tmp1 := strings.Split(path, "/")
	tmp2 := tmp1[0 : len(tmp1)-1] // cut filename
	dir := strings.Join(tmp2, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0o700) // Create all parent directories
	}
}
