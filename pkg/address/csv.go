package address

import (
	"bufio"
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
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
func (c *CSVRepository) CreateFilePath(accountType account.AccountType, keyStatus uint8) string {

	// ./data/pubkey/client_1534744535097796209.csv
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%d_%s.csv", c.filePath, accountType.String(), keyStatus, ts)
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
