package eth

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	pkgrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// Compile-time check that Ethereum implements Ethereumer interface
var _ apieth.Ethereumer = (*Ethereum)(nil)

// Compile-time check that Ethereum implements EtherTxMonitor interface
var _ apieth.EtherTxMonitor = (*Ethereum)(nil)

// Compile-time check that Ethereum implements TxSigner interface
var _ apieth.TxSigner = (*Ethereum)(nil)

// Ethereum includes client to call JSON-RPC
type Ethereum struct {
	ethClient    *ethclient.Client
	pkgrpc       pkgrpc.ETHRPC
	chainConf    *chaincfg.Params
	coinTypeCode domainCoin.CoinTypeCode
	uuidHandler  uuid.UUIDHandler
	conf         *config.Ethereum
	netID        uint16
	version      string
	keyDir       string
	clientType   ClientVersion
}

// NewEthereum creates ethereum object
func NewEthereum(
	ctx context.Context,
	ethClient *ethclient.Client,
	rpcClient *ethrpc.Client,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *config.Ethereum,
	uuidHandler uuid.UUIDHandler,
) (*Ethereum, error) {
	eth := &Ethereum{
		ethClient:    ethClient,
		pkgrpc:       pkgrpc.NewRPCClient(rpcClient),
		coinTypeCode: coinTypeCode,
		uuidHandler:  uuidHandler,
		conf:         conf,
		keyDir:       conf.KeyDirName,
	}

	// key dir
	if eth.keyDir == "" {
		dirName, err := eth.AdminDataDir(ctx)
		if err != nil {
			// Anvil doesn't support admin_datadir RPC - use default
			logger.Warn("admin_datadir RPC not supported, using default keydir")
			eth.keyDir = "./data/keystore"
		} else {
			eth.keyDir = dirName + "/keystore"
		}
	}
	logger.Debug("eth.keyDir", "eth.keyDir", eth.keyDir)

	// get NetID
	netID, err := eth.NetVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.NetVersion(): %w", err)
	}
	eth.netID = netID

	if netID == 1 {
		eth.chainConf = &chaincfg.MainNetParams
	} else {
		eth.chainConf = &chaincfg.TestNet3Params
	}

	// get client version
	clientVer, err := eth.ClientVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.ClientVersion(): %w", err)
	}
	eth.version = clientVer

	eth.clientType = DetectClientType(clientVer)
	logger.Debug("detected client type", "clientType", eth.clientType)

	// check sync progress
	res, isSyncing, err := eth.Syncing(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.Syncing(): %w", err)
	}
	if isSyncing {
		logger.Warn("sync is not completed yet")
	}
	if res != nil {
		logger.Info("still syncing",
			"startingBlock", res.StartingBlock,
			"currentBlock", res.CurrentBlock,
			"highestBlock", res.HighestBlock,
		)
	}

	// check network connections
	isListening, err := eth.NetListening(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call eth.NetListening(): %w", err)
	}
	if !isListening {
		logger.Warn("network is not working")
	}

	return eth, nil
}

// Close disconnect to server
func (e *Ethereum) Close() {
	e.pkgrpc.Close()
}

// CoinTypeCode returns coinTypeCode
func (e *Ethereum) CoinTypeCode() domainCoin.CoinTypeCode {
	return e.coinTypeCode
}

// GetChainConf returns chain conf
func (e *Ethereum) GetChainConf() *chaincfg.Params {
	return e.chainConf
}
