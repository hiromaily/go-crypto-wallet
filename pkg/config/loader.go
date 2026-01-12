package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Environment variable prefix for wallet configuration.
// Example: WALLET_ADDRESS_TYPE, WALLET_BITCOIN_HOST
const envPrefix = "WALLET"

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

	// Enable environment variable override
	// Priority: Environment Variables > Config File > Default Values
	//
	// Environment variable naming convention:
	//   - Prefix: WALLET_
	//   - Nested keys use underscore: WALLET_BITCOIN_HOST
	//   - All uppercase: WALLET_ADDRESS_TYPE
	//
	// Examples:
	//   WALLET_ADDRESS_TYPE=legacy      -> address_type: "legacy"
	//   WALLET_BITCOIN_HOST=127.0.0.1   -> bitcoin.host: "127.0.0.1"
	//   WALLET_MYSQL_HOST=localhost     -> mysql.host: "localhost"
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// Allow running without config file if env vars provide all required config
		// This is useful for containerized deployments (12-factor apps)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error occurred
			return fmt.Errorf("can't read config file. %s: %w", path, err)
		}
		// Config file not found; continue with env vars only
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
