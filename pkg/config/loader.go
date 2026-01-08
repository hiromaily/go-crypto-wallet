package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// loadConfig loads a configuration file and unmarshals it into the target.
//
// This function automatically detects the configuration format based on file extension:
//   - .toml -> TOML format
//   - .yaml, .yml -> YAML format
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
	// Validate file extension and determine config type
	ext := strings.ToLower(filepath.Ext(path))
	var configType string
	switch ext {
	case ".toml":
		configType = "toml"
	case ".yaml", ".yml":
		configType = "yaml"
	default:
		return fmt.Errorf("unsupported config file extension: %s (supported: .toml, .yaml, .yml)", ext)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(configType)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("can't read config file. %s: %w", path, err)
	}

	if err := v.Unmarshal(target); err != nil {
		return fmt.Errorf("fail to unmarshal config: %w", err)
	}

	return nil
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
