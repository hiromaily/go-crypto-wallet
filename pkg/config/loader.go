package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ConfigFormat represents supported configuration file formats.
type ConfigFormat string

const (
	// FormatTOML represents TOML configuration format.
	FormatTOML ConfigFormat = "toml"
	// FormatYAML represents YAML configuration format.
	FormatYAML ConfigFormat = "yaml"
)

// ConfigLoader defines the interface for loading configuration files.
//
// This interface abstracts the configuration loading mechanism to support
// multiple configuration formats (TOML, YAML) without code changes.
// The Strategy pattern is used to allow format-specific implementations.
type ConfigLoader interface {
	// Load reads the configuration file and unmarshals it into the target.
	// The target must be a pointer to a struct.
	Load(path string, target any) error

	// Format returns the configuration format this loader handles.
	Format() ConfigFormat
}

// tomlLoader implements ConfigLoader for TOML files.
type tomlLoader struct{}

// Load reads a TOML configuration file and unmarshals it into the target.
func (*tomlLoader) Load(path string, target any) error {
	return loadWithViper(path, string(FormatTOML), target)
}

// Format returns the TOML format identifier.
func (*tomlLoader) Format() ConfigFormat {
	return FormatTOML
}

// yamlLoader implements ConfigLoader for YAML files.
type yamlLoader struct{}

// Load reads a YAML configuration file and unmarshals it into the target.
func (*yamlLoader) Load(path string, target any) error {
	return loadWithViper(path, string(FormatYAML), target)
}

// Format returns the YAML format identifier.
func (*yamlLoader) Format() ConfigFormat {
	return FormatYAML
}

// loadWithViper is a shared helper that uses viper to load and unmarshal config.
//
// This function handles the common pattern used by all config loaders:
//  1. Creating a viper instance
//  2. Setting the config file path and type
//  3. Reading the config file
//  4. Unmarshaling into the target structure
//
// The target must be a pointer to a struct.
func loadWithViper(path, format string, target any) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(format)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("can't read config file. %s: %w", path, err)
	}

	if err := v.Unmarshal(target); err != nil {
		return fmt.Errorf("fail to unmarshal config: %w", err)
	}

	return nil
}

// NewConfigLoader returns the appropriate ConfigLoader based on file extension.
//
// This factory function implements the Strategy pattern by selecting the
// appropriate loader implementation based on the file extension:
//   - .toml -> tomlLoader
//   - .yaml, .yml -> yamlLoader
//
// Returns an error if the file extension is not supported.
func NewConfigLoader(path string) (ConfigLoader, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".toml":
		return &tomlLoader{}, nil
	case ".yaml", ".yml":
		return &yamlLoader{}, nil
	default:
		return nil, fmt.Errorf("unsupported config file extension: %s (supported: .toml, .yaml, .yml)", ext)
	}
}

// loadConfig loads a configuration file using the appropriate loader.
//
// This is the main entry point for loading configuration files.
// It automatically detects the format based on file extension and
// uses the appropriate loader implementation.
//
// The target must be a pointer to a struct.
//
// Example usage:
//
//	var conf WalletRoot
//	if err := loadConfig("config/wallet.toml", &conf); err != nil {
//	    return err
//	}
//
//	var conf AccountRoot
//	if err := loadConfig("config/account.yaml", &conf); err != nil {
//	    return err
//	}
func loadConfig(path string, target any) error {
	loader, err := NewConfigLoader(path)
	if err != nil {
		return err
	}
	return loader.Load(path, target)
}

// loadTOML loads a TOML configuration file and unmarshals it into the provided target.
//
// This is a generic loader function used by both wallet and account configuration loaders.
// It handles the common pattern of:
//  1. Creating a viper instance
//  2. Setting the config file path and type
//  3. Reading the config file
//  4. Unmarshaling into the target structure
//
// The target must be a pointer to a struct that can be unmarshaled from TOML.
//
// Deprecated: Use loadConfig instead for format-agnostic loading.
// This function is kept for backward compatibility.
func loadTOML[T any](path string) (*T, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("can't read config file. %s: %w", path, err)
	}

	var config T
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("fail to unmarshal config: %w", err)
	}

	return &config, nil
}
