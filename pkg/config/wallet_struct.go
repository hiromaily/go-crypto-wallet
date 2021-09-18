package config

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TODO:
// - use https://github.com/spf13/viper

// WalletRoot wallet root config
type WalletRoot struct {
	AddressType  address.AddrType  `toml:"address_type" validate:"oneof=p2sh-segwit bech32 bch-cashaddr"`
	CoinTypeCode coin.CoinTypeCode `toml:"coin_type"`
	Bitcoin      Bitcoin           `toml:"bitcoin"`
	Ethereum     Ethereum          `toml:"ethereum"`
	Ripple       Ripple            `toml:"ripple"`
	Logger       Logger            `toml:"logger"`
	Tracer       Tracer            `toml:"tracer"`
	MySQL        MySQL             `toml:"mysql"`
	FilePath     FilePath          `toml:"file_path"`
}

// Bitcoin Bitcoin information
type Bitcoin struct {
	Host        string `toml:"host" validate:"required"`
	User        string `toml:"user" validate:"required"`
	Pass        string `toml:"pass" validate:"required"`
	PostMode    bool   `toml:"http_post_mode"`
	DisableTLS  bool   `toml:"disable_tls"`
	NetworkType string `toml:"network_type" validate:"oneof=mainnet testnet3 regtest"`

	Block BitcoinBlock `toml:"block"`
	Fee   BitcoinFee   `toml:"fee"`
}

// BitcoinBlock block information of Bitcoin
// FIXME: keygen/signature wallet doesn't have this value
//  so validation can not be used
type BitcoinBlock struct {
	ConfirmationNum uint64 `toml:"confirmation_num"`
}

// BitcoinFee range of adjustment calculated fee when sending coin
type BitcoinFee struct {
	AdjustmentMin float64 `toml:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max"`
}

// Ethereum information
type Ethereum struct {
	Host            string                    `toml:"host" validate:"required"`
	IPCPath         string                    `toml:"ipc_path"`
	Port            int                       `toml:"port" validate:"required"`
	DisableTLS      bool                      `toml:"disable_tls"`
	NetworkType     string                    `toml:"network_type" validate:"oneof=mainnet goerli rinkeby ropsten"`
	KeyDirName      string                    `toml:"keydir"`
	ConfirmationNum uint64                    `toml:"confirmation_num"`
	ERC20Token      coin.ERC20Token           `toml:"erc20_token"`
	ERC20s          map[coin.ERC20Token]ERC20 `toml:"erc20s"`
}

// ERC20 information
type ERC20 struct {
	Symbol          string `toml:"symbol"`
	Name            string `toml:"name"`
	ContractAddress string `toml:"contract_address"`
	MasterAddress   string `toml:"master_address"`
	Decimals        int    `toml:"decimals"`
}

// Ripple information
type Ripple struct {
	WebsocketPublicURL string    `toml:"websocket_public_url"`
	WebsocketAdminURL  string    `toml:"websocket_admin_url"`
	NetworkType        string    `toml:"network_type" validate:"oneof=mainnet testnet devnet"`
	API                RippleAPI `toml:"api"`
}

// RippleAPI is ripple-lib server info
type RippleAPI struct {
	URL      string       `toml:"url"`
	IsSecure bool         `toml:"is_secure"`
	TxData   RippleTxData `toml:"transaction"`
}

// RippleTxData is used for api command to send coin
type RippleTxData struct {
	Account string `toml:"sender_account"`
	Secret  string `toml:"sender_secret"`
}

// Logger logger info
type Logger struct {
	Service  string `toml:"service" validate:"required"`
	Env      string `toml:"env" validate:"oneof=dev prod custom"`
	Level    string `toml:"level" validate:"required"`
	IsLogger bool   `toml:"is_logger"`
}

// Tracer is open tracing
type Tracer struct {
	Type    string       `toml:"type" validate:"oneof=none jaeger datadog"`
	Jaeger  TracerDetail `toml:"jaeger"`
	Datadog TracerDetail `toml:"datadog"`
}

// TracerDetail includes specific service config
type TracerDetail struct {
	ServiceName         string  `toml:"service_name"`
	CollectorEndpoint   string  `toml:"collector_endpoint"`
	SamplingProbability float64 `toml:"sampling_probability"`
	IsDebug             bool    `toml:"is_debug"`
}

// MySQL MySQL info
type MySQL struct {
	Host  string `toml:"host" validate:"required"`
	DB    string `toml:"dbname" validate:"required"`
	User  string `toml:"user" validate:"required"`
	Pass  string `toml:"pass" validate:"required"`
	Debug bool   `toml:"debug"`
}

// FilePath if file path group
type FilePath struct {
	Tx         string `toml:"tx" validate:"required"`
	Address    string `toml:"address" validate:"required"`
	FullPubKey string `toml:"full_pubkey" validate:"required"`
}

// PubKeyFile saved pubKey file path which is used when import/export file
type PubKeyFile struct {
	BasePath string `toml:"base_path" validate:"required"`
}

// AddressFile saved address file path which is used when import/export file
type AddressFile struct {
	BasePath string `toml:"base_path" validate:"required"`
}
