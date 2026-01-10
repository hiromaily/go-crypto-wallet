package descriptor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// DescriptorFileRepository handles reading and writing descriptor files.
//
// Supports two formats:
//  1. Plain text: One descriptor per line
//  2. Bitcoin Core JSON: JSON array of ImportDescriptorsRequest objects
//
// The repository automatically detects the format when reading and can write
// in either format based on the method called.
type DescriptorFileRepository struct {
	basePath string
}

// NewDescriptorFileRepository creates a new descriptor file repository.
//
// Parameters:
//   - basePath: Base directory for descriptor files
//
// Returns:
//   - DescriptorFileRepository instance
func NewDescriptorFileRepository(basePath string) *DescriptorFileRepository {
	return &DescriptorFileRepository{
		basePath: basePath,
	}
}

// WriteDescriptorsPlainText writes descriptors to a plain text file.
//
// Format: One descriptor per line, optionally with checksum.
//
// Parameters:
//   - filePath: Output file path
//   - descriptors: List of descriptor objects to write
//   - withChecksum: If true, includes checksums in output
//
// Returns:
//   - Error if write fails
//
// Example output:
//
//	sh(wpkh([a1b2c3d4/49'/0'/0']xpub.../0/*))#checksum
//	sh(wpkh([a1b2c3d4/49'/0'/1']xpub.../0/*))#checksum
func (r *DescriptorFileRepository) WriteDescriptorsPlainText(
	filePath string,
	descriptors []*domainWallet.Descriptor,
	withChecksum bool,
) error {
	if len(descriptors) == 0 {
		return errors.New("no descriptors to write")
	}

	// Build output content
	var lines []string
	for _, desc := range descriptors {
		line := desc.Script
		if withChecksum && desc.Checksum != "" {
			line += "#" + desc.Checksum
		}
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n") + "\n"

	// Write file
	return r.writeFile(filePath, []byte(content))
}

// BitcoinCoreImportDescriptor represents a Bitcoin Core importdescriptors request item.
//
// This matches the format expected by Bitcoin Core's importdescriptors RPC.
//
// Reference: Bitcoin Core RPC documentation
type BitcoinCoreImportDescriptor struct {
	// Descriptor is the output descriptor string WITH checksum
	Descriptor string `json:"desc"`

	// Timestamp indicates when to start scanning for transactions
	// "now" = skip rescanning, 0 = scan from genesis
	Timestamp any `json:"timestamp"`

	// Active indicates whether to track outputs for spending
	Active bool `json:"active"`

	// Range specifies the range of addresses to derive
	// Format: [start, end] or integer for 0-N
	Range any `json:"range,omitempty"`

	// NextIndex specifies the next unused address index
	// Only used for ranged descriptors
	NextIndex *int `json:"next_index,omitempty"`

	// Label is an optional label for the imported addresses
	Label string `json:"label,omitempty"`

	// Internal indicates if this is for change addresses
	Internal bool `json:"internal,omitempty"`

	// Watchonly indicates this is a watch-only import (no private keys)
	Watchonly bool `json:"watchonly,omitempty"`
}

// WriteDescriptorsBitcoinCoreJSON writes descriptors in Bitcoin Core JSON format.
//
// This format can be directly used with Bitcoin Core's importdescriptors RPC.
//
// Parameters:
//   - filePath: Output file path
//   - descriptors: List of descriptor objects
//   - timestamp: Rescan timestamp ("now", 0, or unix timestamp)
//   - label: Optional label for imported addresses
//   - active: Whether to track outputs
//   - watchonly: Whether this is watch-only (no private keys)
//
// Returns:
//   - Error if write fails
//
// Example output:
//
//	[
//	  {
//	    "desc": "sh(wpkh([a1b2c3d4/49'/0'/0']xpub.../0/*))#checksum",
//	    "timestamp": "now",
//	    "active": true,
//	    "watchonly": true,
//	    "label": "deposit"
//	  }
//	]
func (r *DescriptorFileRepository) WriteDescriptorsBitcoinCoreJSON(
	filePath string,
	descriptors []*domainWallet.Descriptor,
	timestamp any,
	label string,
	active bool,
	watchonly bool,
) error {
	if len(descriptors) == 0 {
		return errors.New("no descriptors to write")
	}

	// Build Bitcoin Core import requests
	var requests []BitcoinCoreImportDescriptor
	for _, desc := range descriptors {
		// Build descriptor string with checksum
		descriptorStr := desc.Script
		if desc.Checksum != "" {
			descriptorStr += "#" + desc.Checksum
		}

		req := BitcoinCoreImportDescriptor{
			Descriptor: descriptorStr,
			Timestamp:  timestamp,
			Active:     active,
			Watchonly:  watchonly,
		}

		if label != "" {
			req.Label = label
		}

		requests = append(requests, req)
	}

	// Marshal to JSON with pretty printing
	jsonData, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal descriptors to JSON: %w", err)
	}

	// Write file
	return r.writeFile(filePath, jsonData)
}

// ReadDescriptorsPlainText reads descriptors from a plain text file.
//
// Format: One descriptor per line, checksums are automatically detected and parsed.
//
// Parameters:
//   - filePath: Input file path
//
// Returns:
//   - List of descriptor strings (includes checksum if present)
//   - Error if read fails
//
//nolint:revive // Receiver kept for interface compatibility
func (r *DescriptorFileRepository) ReadDescriptorsPlainText(filePath string) ([]string, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open descriptor file %s: %w", filePath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	var descriptors []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		descriptors = append(descriptors, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading descriptor file: %w", err)
	}

	if len(descriptors) == 0 {
		return nil, errors.New("no descriptors found in file")
	}

	return descriptors, nil
}

// ReadDescriptorsBitcoinCoreJSON reads descriptors from Bitcoin Core JSON format.
//
// Parameters:
//   - filePath: Input file path
//
// Returns:
//   - List of BitcoinCoreImportDescriptor objects
//   - Error if read or parse fails
//
//nolint:revive // Receiver kept for interface compatibility
func (r *DescriptorFileRepository) ReadDescriptorsBitcoinCoreJSON(
	filePath string,
) ([]BitcoinCoreImportDescriptor, error) {
	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to read descriptor file %s: %w", filePath, err)
	}

	var descriptors []BitcoinCoreImportDescriptor
	if err := json.Unmarshal(data, &descriptors); err != nil {
		return nil, fmt.Errorf("failed to parse JSON descriptor file: %w", err)
	}

	if len(descriptors) == 0 {
		return nil, errors.New("no descriptors found in JSON file")
	}

	return descriptors, nil
}

// AutoDetectAndRead reads a descriptor file and auto-detects the format.
//
// The format is detected by attempting to parse as JSON first.
// If JSON parsing fails, treats the file as plain text.
//
// Parameters:
//   - filePath: Input file path
//
// Returns:
//   - List of descriptor strings
//   - Error if read fails
func (r *DescriptorFileRepository) AutoDetectAndRead(filePath string) ([]string, error) {
	// Try JSON format first
	jsonDescriptors, err := r.ReadDescriptorsBitcoinCoreJSON(filePath)
	if err == nil {
		// JSON format detected - extract descriptor strings
		descriptors := make([]string, len(jsonDescriptors))
		for i, jd := range jsonDescriptors {
			descriptors[i] = jd.Descriptor
		}
		return descriptors, nil
	}

	// Fall back to plain text format
	return r.ReadDescriptorsPlainText(filePath)
}

// CreateFilePath generates a file path for descriptor export.
//
// Parameters:
//   - accountName: Account name (e.g., "deposit", "payment")
//   - extension: File extension (e.g., ".txt", ".json")
//
// Returns:
//   - Complete file path
func (r *DescriptorFileRepository) CreateFilePath(accountName string, extension string) string {
	return filepath.Join(r.basePath, fmt.Sprintf("%s_descriptors%s", accountName, extension))
}

// WriteFile writes raw data to a file (implements DescriptorFileWriter interface).
//
// This is a generic write method that accepts pre-formatted data.
// For structured descriptor writing, use WriteDescriptorsPlainText or
// WriteDescriptorsBitcoinCoreJSON instead.
//
// Parameters:
//   - path: Output file path
//   - data: Raw data to write
//
// Returns:
//   - Error if write fails
func (r *DescriptorFileRepository) WriteFile(path string, data []byte) error {
	return r.writeFile(path, data)
}

// writeFile is a helper method to write data to a file with proper permissions.
//
//nolint:revive // Receiver kept for consistency
func (r *DescriptorFileRepository) writeFile(filePath string, data []byte) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	// Clean and resolve path
	cleanPath := filepath.Clean(filePath)

	// Create directory if it doesn't exist
	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file with restricted permissions (owner read/write only)
	if err := os.WriteFile(cleanPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write file %s: %w", cleanPath, err)
	}

	return nil
}
