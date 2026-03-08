package mpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	portfile "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
)

// Compile-time check: MPCFileRepository must satisfy MPCFileRepositorier.
var _ portfile.MPCFileRepositorier = (*MPCFileRepository)(nil)

// MPCFileRepository implements MPCFileRepositorier for ETHMPCTransactionFile JSON I/O.
type MPCFileRepository struct {
	filePath string // base directory (with trailing slash)
}

// NewMPCFileRepository creates a new MPCFileRepository with the given base directory.
func NewMPCFileRepository(filePath string) *MPCFileRepository {
	return &MPCFileRepository{filePath: filePath}
}

// CreateMPCFilePath returns the canonical path for an MPC transaction JSON file.
//
//	Unsigned: {filePath}{actionType}_mpc_{uuid}.json
//	Signed:   {filePath}{actionType}_mpc_{uuid}_signed.json
func (r *MPCFileRepository) CreateMPCFilePath(actionType, uuid string, isSigned bool) string {
	if isSigned {
		return fmt.Sprintf("%s%s_mpc_%s_signed.json", r.filePath, actionType, uuid)
	}
	return fmt.Sprintf("%s%s_mpc_%s.json", r.filePath, actionType, uuid)
}

// WriteETHMPCJSONFile validates file, derives the canonical path via CreateMPCFilePath,
// and writes indented JSON. Returns the written file path on success.
func (r *MPCFileRepository) WriteETHMPCJSONFile(file *dtoeth.ETHMPCTransactionFile, isSigned bool) (string, error) {
	if file == nil {
		return "", errors.New("ETHMPCTransactionFile must not be nil")
	}

	if err := file.Validate(); err != nil {
		return "", fmt.Errorf("invalid ETH MPC transaction file: %w", err)
	}

	path := r.CreateMPCFilePath(file.ActionType, file.UUID, isSigned)
	cleanPath := filepath.Clean(path)

	if err := r.checkPathTraversal(cleanPath); err != nil {
		return "", err
	}

	if err := mkdirForFile(cleanPath); err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal ETH MPC transaction: %w", err)
	}

	// 0o600: MPC transaction files contain key-related hashes and signed tx bytes.
	if err = os.WriteFile(cleanPath, jsonData, 0o600); err != nil {
		return "", fmt.Errorf("failed to write ETH MPC transaction file %s: %w", cleanPath, err)
	}

	return cleanPath, nil
}

// ReadETHMPCJSONFile reads, parses, and validates the MPC transaction file at filePath.
func (r *MPCFileRepository) ReadETHMPCJSONFile(filePath string) (*dtoeth.ETHMPCTransactionFile, error) {
	if !strings.HasSuffix(strings.ToLower(filePath), ".json") {
		return nil, fmt.Errorf("invalid MPC file extension: %s (expected .json)", filePath)
	}

	cleanPath := filepath.Clean(filePath)
	if err := r.checkPathTraversal(cleanPath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ETH MPC transaction file %s: %w", cleanPath, err)
	}

	var f dtoeth.ETHMPCTransactionFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("failed to parse ETH MPC transaction file %s: %w", cleanPath, err)
	}

	if err := f.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ETH MPC transaction file data in %s: %w", cleanPath, err)
	}

	return &f, nil
}

// checkPathTraversal validates that cleanPath stays within r.filePath.
// The check is skipped when r.filePath is empty (e.g., in unit tests).
func (r *MPCFileRepository) checkPathTraversal(cleanPath string) error {
	if r.filePath == "" {
		return nil
	}
	baseDir := filepath.Clean(r.filePath)
	rel, err := filepath.Rel(baseDir, cleanPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path traversal attempt detected: %s", cleanPath)
	}
	return nil
}

// mkdirForFile creates all parent directories for the given file path.
func mkdirForFile(cleanPath string) error {
	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}
