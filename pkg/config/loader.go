package config

import (
	"fmt"

	"github.com/spf13/viper"
)

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
