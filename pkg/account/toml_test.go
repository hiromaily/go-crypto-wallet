package account

import (
	"log"
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
