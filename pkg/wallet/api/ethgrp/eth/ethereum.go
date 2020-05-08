package eth

import (
	"context"

	"github.com/btcsuite/btcd/chaincfg"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// Ethereum includes client to call JSON-RPC
type Ethereum struct {
	client       *ethrpc.Client
	chainConf    *chaincfg.Params
	coinTypeCode coin.CoinTypeCode //eth
	logger       *zap.Logger
	ctx          context.Context
	netID        uint16
	version      string
	keyDir      string
}

// NewEthereum creates ethereum object
func NewEthereum(
	ctx context.Context,
	client *ethrpc.Client,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Ethereum,
	logger *zap.Logger) (*Ethereum, error) {

	eth := &Ethereum{
		client:       client,
		coinTypeCode: coinTypeCode,
		logger:       logger,
		ctx:          ctx,
		keyDir:      conf.Geth.KeyDirName,
	}

	// get NetID
	netID, err := eth.NetVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NetVersion()")
	}
	eth.netID = netID

	if netID == 1 {
		eth.chainConf = &chaincfg.MainNetParams
	} else {
		eth.chainConf = &chaincfg.TestNet3Params
	}

	// get version
	clientVer, err := eth.ClientVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.ClientVersion()")
	}
	eth.version = clientVer

	// check sync progress
	res, isSyncing, err := eth.Syncing()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.Syncing()")
	}
	if isSyncing {
		logger.Warn("sync is not completed yet")
	}
	if res != nil {
		logger.Info("still syncing",
			zap.Int64("knownStates", res.KnownStates),
			zap.Int64("pulledStates", res.PulledStates),
			zap.Int64("startingBlock", res.StartingBlock),
			zap.Int64("currentBlock", res.CurrentBlock),
			zap.Int64("highestBlock", res.HighestBlock),
		)
	}

	return eth, nil
}

// Close disconnect to server
func (e *Ethereum) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

// GetChainConf returns chain conf
func (e *Ethereum) GetChainConf() *chaincfg.Params {
	return e.chainConf
}
