package coldwallet

import (
	"go.uber.org/zap"
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
	"github.com/opentracing/opentracing-go"
)

// common interface for keygen/signature
type Coldwalleter interface {
	SignatureFromFile(filePath string) (string, bool, string, error)
	GenerateSeed() ([]byte, error)
	GenerateAccountKey(accountType account.AccountType, coinType enum.CoinType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivateKey(accountType account.AccountType) error
}

// ColdWallet coldwallet for keygen/signature object
type ColdWallet struct {
	btc      api.Bitcoiner
	logger   *zap.Logger
	tracer   opentracing.Tracer
	keyStorager rdb.KeygenStorager
	signStorager rdb.SignatureStorager
	wtype    types.WalletType
}

func NewColdWalet(
	btc api.Bitcoiner,
	logger *zap.Logger,
	tracer opentracing.Tracer,
	keyStorager rdb.KeygenStorager,
	signStorager rdb.SignatureStorager,
	wtype types.WalletType) *ColdWallet {

	return &ColdWallet{
		btc:      btc,
		logger:   logger,
		tracer:   tracer,
		keyStorager: keyStorager,
		signStorager: signStorager,
		wtype:    wtype,
	}
}
