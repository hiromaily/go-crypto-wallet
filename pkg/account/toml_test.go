package account

import (
	"log"
	"os"
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/config/file"
)

// TestNewAccount is test for NewAccount
func TestNewAccount(t *testing.T) {
	// t.SkipNow()
	confPath := file.GetConfigFilePath("account.toml")
	conf, err := NewAccount(confPath)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	grok.Value(conf)
}

// TestNewAccountWithViper tests account loading using viper
func TestNewAccountWithViper(t *testing.T) {
	confPath := file.GetConfigFilePath("account.toml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", confPath)
	}

	conf, err := NewAccount(confPath)
	if err != nil {
		t.Fatalf("NewAccount() error = %v", err)
	}

	if conf == nil {
		t.Fatal("NewAccount() returned nil config")
	}

	// Verify that viper properly loaded the TOML file
	if len(conf.Types) == 0 {
		t.Error("Account types should not be empty")
	}
	if conf.DepositReceiver == "" {
		t.Error("DepositReceiver should not be empty")
	}
	if conf.PaymentSender == "" {
		t.Error("PaymentSender should not be empty")
	}
	if len(conf.Multisigs) == 0 {
		t.Error("Multisigs should not be empty")
	}

	// Verify multisig structure loaded correctly
	for i, ms := range conf.Multisigs {
		if ms.Type == "" {
			t.Errorf("Multisig[%d].Type should not be empty", i)
		}
		if ms.Required == 0 {
			t.Errorf("Multisig[%d].Required should not be zero", i)
		}
		if len(ms.AuthUsers) == 0 {
			t.Errorf("Multisig[%d].AuthUsers should not be empty", i)
		}
	}
}

// TestLoadAccount tests the loadAccount function
func TestLoadAccount(t *testing.T) {
	confPath := file.GetConfigFilePath("account.toml")

	// Skip if config file doesn't exist
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		t.Skipf("Config file not found: %s", confPath)
	}

	conf, err := loadAccount(confPath)
	if err != nil {
		t.Fatalf("loadAccount() error = %v", err)
	}

	if conf == nil {
		t.Fatal("loadAccount() returned nil config")
	}
}
