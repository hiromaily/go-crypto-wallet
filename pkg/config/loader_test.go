package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig tests the loadConfig function with various file formats.
// This verifies that both TOML and YAML files can be loaded automatically.
func TestLoadConfig(t *testing.T) {
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

	// Create YML test file
	ymlPath := filepath.Join(tmpDir, "test.yml")
	err = os.WriteFile(ymlPath, []byte(yamlContent), 0o600)
	require.NoError(t, err, "Failed to create test YML file")

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
			name:    "Load YML file",
			path:    ymlPath,
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
				require.Error(t, err, "loadConfig() should return error")
				return
			}

			require.NoError(t, err, "loadConfig() should not return error")
			assert.NotEmpty(t, conf.Types, "Types should not be empty")
			assert.Equal(t, "deposit", string(conf.DepositReceiver), "DepositReceiver mismatch")
			assert.Equal(t, "payment", string(conf.PaymentSender), "PaymentSender mismatch")
		})
	}
}

// TestLoadConfigWithUnsupportedExtension tests loadConfig with unsupported file extensions.
func TestLoadConfigWithUnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "JSON extension",
			path: filepath.Join(tmpDir, "test.json"),
		},
		{
			name: "XML extension",
			path: filepath.Join(tmpDir, "test.xml"),
		},
		{
			name: "No extension",
			path: filepath.Join(tmpDir, "test"),
		},
		{
			name: "Uppercase unsupported extension",
			path: filepath.Join(tmpDir, "test.JSON"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf AccountRoot
			err := loadConfig(tt.path, &conf)

			require.Error(t, err, "loadConfig() should return error for unsupported extension")
			assert.Contains(t, err.Error(), "unsupported config file extension",
				"Error message should mention unsupported extension")
		})
	}
}

// TestLoadConfigCaseInsensitive tests that file extension detection is case-insensitive.
func TestLoadConfigCaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()

	tomlContent := `
types = ["client", "deposit"]
deposit_receiver = "deposit"
payment_sender = "payment"
`

	yamlContent := `
types:
  - client
  - deposit
deposit_receiver: deposit
payment_sender: payment
`

	tests := []struct {
		name     string
		filename string
		content  string
	}{
		{name: "Uppercase TOML", filename: "test.TOML", content: tomlContent},
		{name: "Mixed case TOML", filename: "test.Toml", content: tomlContent},
		{name: "Uppercase YAML", filename: "test.YAML", content: yamlContent},
		{name: "Mixed case YAML", filename: "test.Yaml", content: yamlContent},
		{name: "Uppercase YML", filename: "test.YML", content: yamlContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.filename)
			err := os.WriteFile(path, []byte(tt.content), 0o600)
			require.NoError(t, err, "Failed to create test file")

			var conf AccountRoot
			err = loadConfig(path, &conf)

			require.NoError(t, err, "loadConfig() should handle case-insensitive extensions")
			assert.NotEmpty(t, conf.Types, "Types should not be empty")
		})
	}
}
