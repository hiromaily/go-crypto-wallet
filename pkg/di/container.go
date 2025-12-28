package di

import (
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/config/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/converter"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// InfrastructureContainer holds instances from pkg/ directory (reusable components)
type InfrastructureContainer struct {
	Config      *config.WalletRoot
	AccountConf *account.AccountRoot
	UUIDHandler uuid.UUIDHandler
	Converter   converter.Converter
}

// NewInfrastructureContainer creates a new infrastructure container with pkg/ components
func NewInfrastructureContainer(
	conf *config.WalletRoot,
	accountConf *account.AccountRoot,
) *InfrastructureContainer {
	return &InfrastructureContainer{
		Config:      conf,
		AccountConf: accountConf,
		UUIDHandler: uuid.NewGoogleUUIDHandler(),
		Converter:   converter.NewConverter(),
	}
}

// SetupLogger initializes the global logger from configuration
func (c *InfrastructureContainer) SetupLogger() {
	logger.SetGlobal(logger.NewSlogFromConfig(
		c.Config.Logger.Env,
		c.Config.Logger.Level,
		c.Config.Logger.Service,
	))
}
