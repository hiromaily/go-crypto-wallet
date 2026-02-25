package address

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// AddressFileRepository is repository to store pubkey as csv file.
// It implements file.AddressFileRepositorier interface defined in
// internal/application/ports/file.
type AddressFileRepository struct {
	filePath string
}

// NewAddressFileRepository returns AddressFileRepository
func NewAddressFileRepository(filePath string) *AddressFileRepository {
	return &AddressFileRepository{
		filePath: filePath,
	}
}

// CreateFilePath create file path for csv file
// Format:
//   - ./data/pubkey/client_1534744535097796209.csv
func (r *AddressFileRepository) CreateFilePath(accountType domainAccount.AccountType) string {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%s.csv", r.filePath, accountType.String(), ts)
}

// ValidateFilePath validate fileName
func (*AddressFileRepository) ValidateFilePath(fileName string, accountType domainAccount.AccountType) error {
	// e.g. ./data/pubkey/deposit/deposit_1586831083436291000.csv
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != accountType.String() {
		return fmt.Errorf("mismatching between accountType [%s] and file prefix [%s]", accountType, tmp[0])
	}
	return nil
}

// WriteXpubLine creates a new file and writes a single xpub CSV line.
// Format: coinTypeCode,accountType,44,extendedPubKey,derivationPath
func (r *AddressFileRepository) WriteXpubLine(
	accountType domainAccount.AccountType,
	coinTypeCode, xpub, derivationPath string,
) (string, error) {
	fileName := r.CreateFilePath(accountType)

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("failed to create xpub file %s: %w", fileName, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close xpub file: %w", cerr)
		}
	}()

	line := fmt.Sprintf("%s,%s,%d,%s,%s\n",
		coinTypeCode, accountType.String(), domainCoin.BIP44Purpose, xpub, derivationPath)
	if _, err = fmt.Fprint(f, line); err != nil {
		return "", fmt.Errorf("failed to write xpub file: %w", err)
	}
	return fileName, nil
}

// ImportAddress import pubkey from csv file
func (*AddressFileRepository) ImportAddress(fileName string) ([]string, error) {
	file, err := os.Open(fileName) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s) error: %s", fileName, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	var pubKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}
