package descriptor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// FileWriter writes descriptor data to disk.
type FileWriter struct{}

// NewFileWriter creates a new descriptor file writer.
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// WriteFile writes descriptor content to the specified path, creating directories as needed.
func (*FileWriter) WriteFile(path string, data []byte) error {
	if path == "" {
		return errors.New("output path cannot be empty")
	}

	cleanPath := filepath.Clean(path)
	dir := filepath.Dir(cleanPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(cleanPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write descriptor file %s: %w", cleanPath, err)
	}

	return nil
}
