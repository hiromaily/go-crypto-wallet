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

	"github.com/hiromaily/go-bitcoin/pkg/account"
)

// CSVRepository to store pubkey as csv file
type CSVRepository struct {
	filePath string
	logger   *zap.Logger
}

// NewCSVRepository
func NewCSVRepository(filePath string, logger *zap.Logger) *CSVRepository {
	return &CSVRepository{
		filePath: filePath,
		logger:   logger,
	}
}

// CreateFilePath create file path for csv file
// Format:
//  - ./data/pubkey/client_1534744535097796209.csv
func (c *CSVRepository) CreateFilePath(accountType account.AccountType, keyStatus uint8) string {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%d_%s.csv", c.filePath, accountType.String(), keyStatus, ts)
}

// ValidateFilePath validate fileName
func (c *CSVRepository) ValidateFilePath(fileName string, accountType account.AccountType) error {
	//e.g. ./data/pubkey/receipt/receipt_1_1586831083436291000.csv
	tmp := strings.Split(strings.Split(fileName, "_")[0], "/")
	if tmp[len(tmp)-1] != accountType.String() {
		return errors.Errorf("mismatching between accountType [%s] and file prefix [%s]", accountType, tmp[0])
	}
	return nil
}

// ImportPubKey import pubkey from csv file
func (c *CSVRepository) ImportPubKey(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Errorf("os.Open(%s) error: %s", fileName, err)
	}
	defer file.Close()

	var pubKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}
