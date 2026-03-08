// Package mpc provides infrastructure implementations for MPC/TSS key shard
// storage and MPC transaction file I/O.
package mpc

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/scrypt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
)

// Compile-time check: MPCShardStore must satisfy MPCKeyShardStorage.
var _ apieth.MPCKeyShardStorage = (*MPCShardStore)(nil)

// File layout: salt[16] || nonce[12] || ciphertext.
const (
	scryptN  = 1 << 15 // CPU/memory cost parameter (N=32768)
	scryptR  = 8       // block size
	scryptP  = 1       // parallelization factor
	saltLen  = 16      // random salt length in bytes
	nonceLen = 12      // AES-GCM standard nonce length in bytes
	keyLen   = 32      // AES-256 key length in bytes

	// minFileLen is the minimum valid encrypted file size:
	// salt(16) + nonce(12) + GCM tag(16) = 44 bytes of overhead, plus at least 1 byte plaintext.
	minFileLen = saltLen + nonceLen + 16
)

// MPCShardStore implements MPCKeyShardStorage using AES-256-GCM encryption
// with a scrypt-derived key. The full plaintext is never written to disk.
type MPCShardStore struct{}

// NewMPCShardStore creates a new MPCShardStore.
func NewMPCShardStore() *MPCShardStore {
	return &MPCShardStore{}
}

// SaveShard encrypts data with AES-256-GCM (key derived via scrypt from passphrase)
// and writes the ciphertext to shardPath atomically using a temp-file + rename.
// Parent directories are created if they do not exist.
func (*MPCShardStore) SaveShard(_ context.Context, shardPath string, passphrase string, data []byte) error {
	salt := make([]byte, saltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("mpc shard: failed to generate salt: %w", err)
	}

	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return err
	}

	gcm, err := newGCM(key)
	if err != nil {
		return err
	}

	nonce := make([]byte, nonceLen)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("mpc shard: failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// File format: salt || nonce || ciphertext
	fileData := make([]byte, 0, saltLen+nonceLen+len(ciphertext))
	fileData = append(fileData, salt...)
	fileData = append(fileData, nonce...)
	fileData = append(fileData, ciphertext...)

	return atomicWrite(shardPath, fileData)
}

// LoadShard reads and decrypts the shard file at shardPath using the given passphrase.
// Returns the original plaintext bytes, or an error if the file is missing,
// too short, or the passphrase is wrong (authentication failure).
func (*MPCShardStore) LoadShard(_ context.Context, shardPath string, passphrase string) ([]byte, error) {
	cleanPath := filepath.Clean(shardPath)

	fileData, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("mpc shard: failed to read shard file %s: %w", cleanPath, err)
	}

	if len(fileData) < minFileLen {
		return nil, fmt.Errorf("mpc shard: file %s is too short to be a valid shard (got %d bytes, need at least %d)",
			cleanPath, len(fileData), minFileLen)
	}

	salt := fileData[:saltLen]
	nonce := fileData[saltLen : saltLen+nonceLen]
	ciphertext := fileData[saltLen+nonceLen:]

	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return nil, err
	}

	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Do not expose raw crypto error — it may hint at partial key material.
		return nil, fmt.Errorf("mpc shard: decryption failed for %s (wrong passphrase or corrupted file)", cleanPath)
	}

	return plaintext, nil
}

// deriveKey uses scrypt to derive a 32-byte AES-256 key from passphrase and salt.
func deriveKey(passphrase string, salt []byte) ([]byte, error) {
	key, err := scrypt.Key([]byte(passphrase), salt, scryptN, scryptR, scryptP, keyLen)
	if err != nil {
		return nil, fmt.Errorf("mpc shard: failed to derive key: %w", err)
	}
	return key, nil
}

// newGCM constructs an AES-GCM AEAD from a 32-byte key.
func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("mpc shard: failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("mpc shard: failed to create GCM: %w", err)
	}
	return gcm, nil
}

// atomicWrite writes data to path atomically by writing a temp file in the same
// directory and renaming it, preventing partial writes from corrupting the shard.
func atomicWrite(shardPath string, data []byte) error {
	cleanPath := filepath.Clean(shardPath)
	dir := filepath.Dir(cleanPath)

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("mpc shard: failed to create directory %s: %w", dir, err)
	}

	tmpFile, err := os.CreateTemp(dir, ".shard-*.tmp")
	if err != nil {
		return fmt.Errorf("mpc shard: failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("mpc shard: failed to write temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("mpc shard: failed to close temp file: %w", err)
	}

	if err := os.Chmod(tmpPath, 0o600); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("mpc shard: failed to set permissions on temp file: %w", err)
	}

	if err := os.Rename(tmpPath, cleanPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("mpc shard: failed to rename temp file to %s: %w", cleanPath, err)
	}

	return nil
}
