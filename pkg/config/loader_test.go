package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewConfigLoader tests the config loader factory function.
// This verifies that the factory correctly resolves loaders based on file extension.
func TestNewConfigLoader(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantFormat ConfigFormat
		wantErr    bool
	}{
		{
			name:       "TOML file extension",
			path:       "config/test.toml",
			wantFormat: FormatTOML,
			wantErr:    false,
		},
		{
			name:       "YAML file extension (.yaml)",
			path:       "config/test.yaml",
			wantFormat: FormatYAML,
			wantErr:    false,
		},
		{
			name:       "YAML file extension (.yml)",
			path:       "config/test.yml",
			wantFormat: FormatYAML,
			wantErr:    false,
		},
		{
			name:       "Uppercase TOML extension",
			path:       "config/TEST.TOML",
			wantFormat: FormatTOML,
			wantErr:    false,
		},
		{
			name:       "Uppercase YAML extension",
			path:       "config/TEST.YAML",
			wantFormat: FormatYAML,
			wantErr:    false,
		},
		{
			name:    "Unsupported extension (.json)",
			path:    "config/test.json",
			wantErr: true,
		},
		{
			name:    "Unsupported extension (.xml)",
			path:    "config/test.xml",
			wantErr: true,
		},
		{
			name:    "No extension",
			path:    "config/test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader, err := NewConfigLoader(tt.path)

			if tt.wantErr {
				require.Error(t, err, "NewConfigLoader() should return error")
				assert.Nil(t, loader, "Loader should be nil on error")
				return
			}

			require.NoError(t, err, "NewConfigLoader() should not return error")
			require.NotNil(t, loader, "Loader should not be nil")
			assert.Equal(t, tt.wantFormat, loader.Format(), "Format mismatch")
		})
	}
}

// TestConfigLoaderLoad tests the Load method of config loaders.
// This verifies that both TOML and YAML loaders can load configuration files.
func TestConfigLoaderLoad(t *testing.T) {
	// Create temporary test config files
	tmpDir := t.TempDir()

	// Create TOML test file
	tomlContent := `
types = ["client", "deposit"]
deposit_receiver = "deposit"
payment_sender = "payment"
`
	tomlPath := filepath.Join(tmpDir, "test.toml")
	err := os.WriteFile(tomlPath, []byte(tomlContent), 0o600)
	require.NoError(t, err, "Failed to create test TOML file")

	// Create YAML test file
	yamlContent := `
types:
  - client
  - deposit
deposit_receiver: deposit
payment_sender: payment
`
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YAML file")

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Load TOML file",
			path:    tomlPath,
			wantErr: false,
		},
		{
			name:    "Load YAML file",
			path:    yamlPath,
			wantErr: false,
		},
		{
			name:    "Load non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.toml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf AccountRoot
			err := loadConfig(tt.path, &conf)

			if tt.wantErr {
				require.Error(t, err, "Load() should return error")
				return
			}

			require.NoError(t, err, "Load() should not return error")
			assert.NotEmpty(t, conf.Types, "Types should not be empty")
			assert.Equal(t, "deposit", string(conf.DepositReceiver), "DepositReceiver mismatch")
			assert.Equal(t, "payment", string(conf.PaymentSender), "PaymentSender mismatch")
		})
	}
}

// TestLoadConfigWithInvalidExtension tests loadConfig with unsupported file extensions.
func TestLoadConfigWithInvalidExtension(t *testing.T) {
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "test.json")

	var conf AccountRoot
	err := loadConfig(invalidPath, &conf)

	require.Error(t, err, "loadConfig() should return error for unsupported extension")
	assert.Contains(t, err.Error(), "unsupported config file extension",
		"Error message should mention unsupported extension")
}

// TestTOMLLoaderImplementation tests tomlLoader implementation.
func TestTOMLLoaderImplementation(t *testing.T) {
	loader := &tomlLoader{}
	assert.Equal(t, FormatTOML, loader.Format(), "tomlLoader should return TOML format")
}

// TestYAMLLoaderImplementation tests yamlLoader implementation.
func TestYAMLLoaderImplementation(t *testing.T) {
	loader := &yamlLoader{}
	assert.Equal(t, FormatYAML, loader.Format(), "yamlLoader should return YAML format")
}
