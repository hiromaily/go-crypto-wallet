package address

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
)

// FileRepositorier is address storage interface
type FileRepositorier interface {
	CreateFilePath(accountType domainAccount.AccountType) string
	ValidateFilePath(fileName string, accountType domainAccount.AccountType) error
	ImportAddress(fileName string) ([]string, error)
}

// FileRepository is repository to store pubkey as csv file
type FileRepository struct {
	filePath string
}

// NewFileRepository returns FileRepository
func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{
		filePath: filePath,
	}
}

// CreateFilePath create file path for csv file
// Format:
//   - ./data/pubkey/client_1534744535097796209.csv
func (r *FileRepository) CreateFilePath(accountType domainAccount.AccountType) string {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%s.csv", r.filePath, accountType.String(), ts)
}

// ValidateFilePath validate fileName
func (*FileRepository) ValidateFilePath(fileName string, accountType domainAccount.AccountType) error {
	// e.g. ./data/pubkey/deposit/deposit_1586831083436291000.csv
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != accountType.String() {
		return fmt.Errorf("mismatching between accountType [%s] and file prefix [%s]", accountType, tmp[0])
	}
	return nil
}

// ImportAddress import pubkey from csv file
func (*FileRepository) ImportAddress(fileName string) ([]string, error) {
	file, err := os.Open(fileName) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s) error: %s", fileName, err)
	}

	defer file.Close()

	var pubKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}
