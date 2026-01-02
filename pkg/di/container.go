package di

import (
	"database/sql"

	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/serializer"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// PkgContainer is the interface for the package container
type PkgContainer interface {
	NewUUIDHandler() uuid.UUIDHandler
	NewMySQLClient() *sql.DB
	NewLogger() logger.Logger
	NewGRPCClient() *grpc.ClientConn
	NewSerializer() serializer.Serializer
}

var _ PkgContainer = (*pkgContainer)(nil)

// pkgContainer holds instances from pkg/ directory (reusable components)
type pkgContainer struct {
	// config
	config      *config.WalletRoot
	accountConf *config.AccountRoot
	// logger
	logger logger.Logger
	// uuid
	uuidHandler uuid.UUIDHandler
	// db
	mysqlClient *sql.DB
	// grpc
	grpcConn *grpc.ClientConn
	// serial
	serializer serializer.Serializer
}

// NewPkgContainer creates a new package container with pkg/ components
func NewPkgContainer(
	conf *config.WalletRoot,
	accountConf *config.AccountRoot,
) *pkgContainer {
	return &pkgContainer{
		config:      conf,
		accountConf: accountConf,
	}
}

// NewLogger creates a new logger
func (c *pkgContainer) NewLogger() logger.Logger {
	if c.logger == nil {
		c.logger = logger.NewSlogFromConfig(
			c.config.Logger.Env,
			c.config.Logger.Level,
			c.config.Logger.Service,
		)
	}
	return c.logger
}

// NewUUIDHandler creates a new UUID handler
func (c *pkgContainer) NewUUIDHandler() uuid.UUIDHandler {
	if c.uuidHandler == nil {
		c.uuidHandler = uuid.NewGoogleUUIDHandler()
	}
	return c.uuidHandler
}

// NewMySQLClient creates a new MySQL client
func (c *pkgContainer) NewMySQLClient() *sql.DB {
	if c.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&c.config.MySQL)
		if err != nil {
			panic(err)
		}
		c.mysqlClient = dbConn
	}
	return c.mysqlClient
}

// NewGRPCClient creates a new gRPC client
func (c *pkgContainer) NewGRPCClient() *grpc.ClientConn {
	if c.grpcConn == nil {
		grpcConn, err := grpc.NewClient(c.config.Ripple.API.URL)
		if err != nil {
			panic(err)
		}
		c.grpcConn = grpcConn
	}
	return c.grpcConn
}

// NewSerializer creates a new serializer
// Default is gob serializer for backward compatibility
func (c *pkgContainer) NewSerializer() serializer.Serializer {
	if c.serializer == nil {
		c.serializer = serializer.NewGobSerializer()
	}
	return c.serializer
}
