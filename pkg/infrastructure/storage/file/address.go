package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
)

// AddressFileRepositorier is address storage interface
type AddressFileRepositorier interface {
	CreateFilePath(accountType domainAccount.AccountType) string
	ValidateFilePath(fileName string, accountType domainAccount.AccountType) error
	ImportAddress(fileName string) ([]string, error)
}

// AddressFileRepository is repository to store pubkey as csv file
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
