package address

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
)

// FileRepositorier is address storage interface
type FileRepositorier interface {
	CreateFilePath(accountType account.AccountType) string
	ValidateFilePath(fileName string, accountType account.AccountType) error
	ImportAddress(fileName string) ([]string, error)
}

// FileRepository is repository to store pubkey as csv file
type FileRepository struct {
	filePath string
	logger   *zap.Logger
}

// NewFileRepository returns FileRepository
func NewFileRepository(filePath string, logger *zap.Logger) *FileRepository {
	return &FileRepository{
		filePath: filePath,
		logger:   logger,
	}
}

// CreateFilePath create file path for csv file
// Format:
//   - ./data/pubkey/client_1534744535097796209.csv
func (r *FileRepository) CreateFilePath(accountType account.AccountType) string {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%s.csv", r.filePath, accountType.String(), ts)
}

// ValidateFilePath validate fileName
func (r *FileRepository) ValidateFilePath(fileName string, accountType account.AccountType) error {
	// e.g. ./data/pubkey/deposit/deposit_1586831083436291000.csv
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != accountType.String() {
		return errors.Errorf("mismatching between accountType [%s] and file prefix [%s]", accountType, tmp[0])
	}
	return nil
}

// ImportAddress import pubkey from csv file
func (r *FileRepository) ImportAddress(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Errorf("os.Open(%s) error: %s", fileName, err)
	}
	// nolint:errcheck
	defer file.Close()

	var pubKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}
