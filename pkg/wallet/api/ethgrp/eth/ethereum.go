package eth

import (
	"context"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// Ethereum includes client to call JSON-RPC
type Ethereum struct {
	ethClient    *ethclient.Client
	rpcClient    *ethrpc.Client
	chainConf    *chaincfg.Params
	coinTypeCode coin.CoinTypeCode
	logger       *zap.Logger
	ctx          context.Context
	netID        uint16
	version      string
	keyDir       string
	isParity     bool
}

// NewEthereum creates ethereum object
func NewEthereum(
	ctx context.Context,
	ethClient *ethclient.Client,
	rpcClient *ethrpc.Client,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Ethereum,
	logger *zap.Logger) (*Ethereum, error) {
	eth := &Ethereum{
		ethClient:    ethClient,
		rpcClient:    rpcClient,
		coinTypeCode: coinTypeCode,
		logger:       logger,
		ctx:          ctx,
		keyDir:       conf.KeyDirName,
	}

	// key dir
	if eth.keyDir == "" {
		dirName, err := eth.AdminDataDir()
		if err != nil {
			return nil, errors.Wrap(err, "fail to call eth.AdminDataDir()")
		}
		eth.keyDir = fmt.Sprintf("%s/keystore", dirName)
	}
	logger.Debug("eth.keyDir", zap.String("eth.keyDir", eth.keyDir))

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

	// get client version
	clientVer, err := eth.ClientVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.ClientVersion()")
	}
	eth.version = clientVer

	eth.isParity = isParity(clientVer)

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
			zap.Int64("startingBlock", res.StartingBlock),
			zap.Int64("currentBlock", res.CurrentBlock),
			zap.Int64("highestBlock", res.HighestBlock),
		)
	}

	// check network connections
	isListening, err := eth.NetListening()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NetListening()")
	}
	if !isListening {
		logger.Warn("network is not working")
	}

	return eth, nil
}

// Close disconnect to server
func (e *Ethereum) Close() {
	if e.rpcClient != nil {
		e.rpcClient.Close()
	}
}

// CoinTypeCode returns coinTypeCode
func (e *Ethereum) CoinTypeCode() coin.CoinTypeCode {
	return e.coinTypeCode
}

// GetChainConf returns chain conf
func (e *Ethereum) GetChainConf() *chaincfg.Params {
	return e.chainConf
}

func isParity(target string) bool {
	return strings.Contains(target, ClientVersionParity.String())
}
