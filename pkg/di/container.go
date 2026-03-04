package di

import (
	"database/sql"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/postgres"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/sqlite"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/serializer"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// PkgContainer is the interface for the package container
type PkgContainer interface {
	NewUUIDHandler() uuid.UUIDHandler
	NewDatabaseClient() *sql.DB
	NewMySQLClient() *sql.DB
	NewPostgresClient() *sql.DB
	NewSQLiteClient() *sql.DB
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
	mysqlClient    *sql.DB
	postgresClient *sql.DB
	sqliteClient   *sql.DB
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
		c.logger = logger.NewSlogLogger(
			c.config.Logger.Format,
			c.config.Logger.Level,
			c.config.Logger.Service,
			"",
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

// NewDatabaseClient creates a database client based on configured type
func (c *pkgContainer) NewDatabaseClient() *sql.DB {
	switch c.config.Database.Type {
	case "mysql":
		return c.NewMySQLClient()
	case "postgres":
		return c.NewPostgresClient()
	case "sqlite":
		return c.NewSQLiteClient()
	default:
		panic("unsupported database type: " + c.config.Database.Type)
	}
}

// NewMySQLClient creates a new MySQL client
func (c *pkgContainer) NewMySQLClient() *sql.DB {
	if c.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&c.config.Database.MySQL)
		if err != nil {
			panic(err)
		}
		c.mysqlClient = dbConn
	}
	return c.mysqlClient
}

// NewPostgresClient creates a new PostgreSQL client
func (c *pkgContainer) NewPostgresClient() *sql.DB {
	if c.postgresClient == nil {
		dbConn, err := postgres.NewPostgres(&c.config.Database.PostgreSQL)
		if err != nil {
			panic(err)
		}
		c.postgresClient = dbConn
	}
	return c.postgresClient
}

// NewSQLiteClient creates a new SQLite client
func (c *pkgContainer) NewSQLiteClient() *sql.DB {
	if c.sqliteClient == nil {
		dbConn, err := sqlite.NewSQLite(&c.config.Database.SQLite)
		if err != nil {
			panic(err)
		}
		c.sqliteClient = dbConn
	}
	return c.sqliteClient
}

// NewGRPCClient creates a new gRPC client
func (c *pkgContainer) NewGRPCClient() *grpc.ClientConn {
	if c.grpcConn == nil {
		var opts []grpc.DialOption
		if !c.config.Ripple.API.IsSecure {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
		grpcConn, err := grpc.NewClient(c.config.Ripple.API.URL, opts...)
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
