package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/cpacia/bchutil"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/pkg/errors"
)

// Bitcoin includes Client to call Json-RPC
type Bitcoin struct {
	client    *rpcclient.Client
	chainConf *chaincfg.Params
	//receipt           KeyInfo
	//payment           KeyInfo
	confirmationBlock int
	feeRange          FeeAdjustmentRate
	version           enum.BTCVersion //179900
}

// KeyInfo 公開鍵アドレスと紐づくアカウント名
//type KeyInfo struct {
//	address    string
//	acountName string
//}

// FeeAdjustmentRate 手数料調整のRange
type FeeAdjustmentRate struct {
	min float64
	max float64
}

// Connection is to local bitcoin core RPC server using HTTP POST mode
//func Connection(host, user, pass string, postMode, tls, isMain bool) (*Bitcoin, error) {
func Connection(conf *toml.BitcoinConf) (*Bitcoin, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Pass,
		HTTPPostMode: conf.PostMode,   // Bitcoin core only supports HTTP POST mode
		DisableTLS:   conf.DisableTLS, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, errors.Errorf("rpcclient.New() error: %s", err)
	}

	bit := Bitcoin{client: client}
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
	logger.Infof("bitcoin server version: %d", netInfo.Version)

	bit.confirmationBlock = conf.Block.ConfirmationNum
	bit.feeRange.max = conf.Fee.AdjustmentMax
	bit.feeRange.min = conf.Fee.AdjustmentMin
	//bit.receipt.address = conf.Stored.Address
	//bit.receipt.acountName = conf.Stored.AccountName
	//bit.payment.address = conf.Payment.Address
	//bit.payment.acountName = conf.Payment.AccountName

	return &bit, err
}

// Close コネクションを切断する
func (b *Bitcoin) Close() {
	b.client.Shutdown()
}

// GetChainConf 接続先であるMainNet/TestNetに応じて必要なconfを返す
func (b *Bitcoin) GetChainConf() *chaincfg.Params {
	return b.chainConf
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

// Version bitcoin coreのバージョンを返す
func (b *Bitcoin) Version() enum.BTCVersion {
	return b.version
}

// OverrideChainParamsByBCH chaincfgをBCH用に上書きする
func (b *Bitcoin) OverrideChainParamsByBCH() {
	switch b.chainConf.Name {
	case chaincfg.TestNet3Params.Name:
		b.chainConf.Net = bchutil.TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		b.chainConf.Net = bchutil.Regtestmagic
	default:
		//chaincfg.MainNetParams.Name
		b.chainConf.Net = bchutil.MainnetMagic
	}
}

// ReceiptAddress 入金用アドレスを返す
//func (b *Bitcoin) ReceiptAddress() string {
//	return b.receipt.address
//}

// ReceiptAccountName 入金用アカウント名を返す
//func (b *Bitcoin) ReceiptAccountName() string {
//	return b.receipt.acountName
//}

// PaymentAddress 支払い用アドレスを返す
//func (b *Bitcoin) PaymentAddress() string {
//	return b.payment.address
//}

// PaymentAccountName 支払い用アカウント名を返す
//func (b *Bitcoin) PaymentAccountName() string {
//	return b.payment.acountName
//}
