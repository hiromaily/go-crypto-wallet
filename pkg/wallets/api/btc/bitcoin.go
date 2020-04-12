package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
)

// Bitcoin includes Client to call Json-RPC
type Bitcoin struct {
	client            *rpcclient.Client
	chainConf         *chaincfg.Params
	confirmationBlock int
	feeRange          FeeAdjustmentRate
	version           ctype.BTCVersion //179900
	coinType          ctype.CoinType   //btc
	logger            *zap.Logger
}

// FeeAdjustmentRate 手数料調整のRange
type FeeAdjustmentRate struct {
	min float64
	max float64
}

// NewBitcoin Bitcoinオブジェクトを返す
func NewBitcoin(client *rpcclient.Client, conf *config.Bitcoin, logger *zap.Logger) (*Bitcoin, error) {
	bit := Bitcoin{
		client: client,
		logger: logger,
	}
	if conf.IsMain {
		bit.chainConf = &chaincfg.MainNetParams
	} else {
		bit.chainConf = &chaincfg.TestNet3Params
	}

	//Bitcoinのバージョンを入れておく
	netInfo, err := bit.GetNetworkInfo()
	if err != nil {
		return nil, errors.Errorf("bit.GetNetworkInfo() error: %s", err)
	}
	bit.version = netInfo.Version
	bit.logger.Info("bitcoin rpc server", zap.Int("version", netInfo.Version.Int()))

	bit.coinType = ctype.BTC

	bit.confirmationBlock = conf.Block.ConfirmationNum
	bit.feeRange.max = conf.Fee.AdjustmentMax
	bit.feeRange.min = conf.Fee.AdjustmentMin

	return &bit, nil
}

// Close コネクションを切断する
func (b *Bitcoin) Close() {
	b.client.Shutdown()
}

// GetChainConf 接続先であるMainNet/TestNetに応じて必要なconfを返す
func (b *Bitcoin) GetChainConf() *chaincfg.Params {
	return b.chainConf
}

// SetChainConf chainConfをセットする
func (b *Bitcoin) SetChainConf(conf *chaincfg.Params) {
	b.chainConf = conf
}

// SetChainConfNet conf.Netをセットする
func (b *Bitcoin) SetChainConfNet(btcNet wire.BitcoinNet) {
	b.chainConf.Net = btcNet
}

// Client clientオブジェクトを返す
func (b *Bitcoin) Client() *rpcclient.Client {
	return b.client
}

// ConfirmationBlock Confirmationに必要なブロック数を返す
func (b *Bitcoin) ConfirmationBlock() int {
	return b.confirmationBlock
}

// FeeRangeMax feeの調整倍率の最大値を返す
func (b *Bitcoin) FeeRangeMax() float64 {
	return b.feeRange.max
}

// FeeRangeMin feeの調整倍率の最小値を返す
func (b *Bitcoin) FeeRangeMin() float64 {
	return b.feeRange.min
}

// SetVersion バージョン情報をセットする
func (b *Bitcoin) SetVersion(ver ctype.BTCVersion) {
	b.version = ver
}

// Version bitcoin coreのバージョンを返す
func (b *Bitcoin) Version() ctype.BTCVersion {
	return b.version
}

// SetCoinType CoinTypeをセットする
func (b *Bitcoin) SetCoinType(coinType ctype.CoinType) {
	b.coinType = coinType
}

// CoinType Bitcoinの種別(btc, bch)を返す
func (b *Bitcoin) CoinType() ctype.CoinType {
	return b.coinType
}
