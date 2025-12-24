package config

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
)

// WalletRoot wallet root config
type WalletRoot struct {
	//nolint:lll,revive
	AddressType  address.AddrType        `toml:"address_type" mapstructure:"address_type" validate:"oneof=p2sh-segwit bech32 bch-cashaddr"`
	CoinTypeCode domainCoin.CoinTypeCode `toml:"coin_type" mapstructure:"coin_type"`
	Bitcoin      Bitcoin                 `toml:"bitcoin" mapstructure:"bitcoin"`
	Ethereum     Ethereum                `toml:"ethereum" mapstructure:"ethereum"`
	Ripple       Ripple                  `toml:"ripple" mapstructure:"ripple"`
	Logger       Logger                  `toml:"logger" mapstructure:"logger"`
	Tracer       Tracer                  `toml:"tracer" mapstructure:"tracer"`
	MySQL        MySQL                   `toml:"mysql" mapstructure:"mysql"`
	FilePath     FilePath                `toml:"file_path" mapstructure:"file_path"`
}

// Bitcoin information
type Bitcoin struct {
	Host       string `toml:"host" mapstructure:"host" validate:"required"`
	User       string `toml:"user" mapstructure:"user" validate:"required"`
	Pass       string `toml:"pass" mapstructure:"pass" validate:"required"`
	PostMode   bool   `toml:"http_post_mode" mapstructure:"http_post_mode"`
	DisableTLS bool   `toml:"disable_tls" mapstructure:"disable_tls"`
	//nolint:lll,revive
	NetworkType string `toml:"network_type" mapstructure:"network_type" validate:"oneof=mainnet testnet3 regtest signet"`

	Block BitcoinBlock `toml:"block" mapstructure:"block"`
	Fee   BitcoinFee   `toml:"fee" mapstructure:"fee"`
}

// BitcoinBlock block information of Bitcoin
// FIXME: keygen/signature wallet doesn't have this value
//
//	so validation can not be used
type BitcoinBlock struct {
	ConfirmationNum uint64 `toml:"confirmation_num" mapstructure:"confirmation_num"`
}

// BitcoinFee range of adjustment calculated fee when sending coin
type BitcoinFee struct {
	AdjustmentMin float64 `toml:"adjustment_min" mapstructure:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max" mapstructure:"adjustment_max"`
}

// Ethereum information
type Ethereum struct {
	Host       string `toml:"host" mapstructure:"host" validate:"required"`
	IPCPath    string `toml:"ipc_path" mapstructure:"ipc_path"`
	Port       int    `toml:"port" mapstructure:"port" validate:"required"`
	DisableTLS bool   `toml:"disable_tls" mapstructure:"disable_tls"`
	//nolint:lll,revive
	NetworkType     string                          `toml:"network_type" mapstructure:"network_type" validate:"oneof=mainnet goerli rinkeby ropsten"`
	KeyDirName      string                          `toml:"keydir" mapstructure:"keydir"`
	ConfirmationNum uint64                          `toml:"confirmation_num" mapstructure:"confirmation_num"`
	ERC20Token      domainCoin.ERC20Token           `toml:"erc20_token" mapstructure:"erc20_token"`
	ERC20s          map[domainCoin.ERC20Token]ERC20 `toml:"erc20s" mapstructure:"erc20s"`
}

// ERC20 information
type ERC20 struct {
	Symbol          string `toml:"symbol" mapstructure:"symbol"`
	Name            string `toml:"name" mapstructure:"name"`
	ContractAddress string `toml:"contract_address" mapstructure:"contract_address"`
	MasterAddress   string `toml:"master_address" mapstructure:"master_address"`
	Decimals        int    `toml:"decimals" mapstructure:"decimals"`
}

// Ripple information
type Ripple struct {
	WebsocketPublicURL string `toml:"websocket_public_url" mapstructure:"websocket_public_url"`
	WebsocketAdminURL  string `toml:"websocket_admin_url" mapstructure:"websocket_admin_url"`
	//nolint:lll
	NetworkType string    `toml:"network_type" mapstructure:"network_type" validate:"oneof=mainnet testnet devnet"`
	API         RippleAPI `toml:"api" mapstructure:"api"`
}

// RippleAPI is ripple-lib server info
type RippleAPI struct {
	URL      string       `toml:"url" mapstructure:"url"`
	IsSecure bool         `toml:"is_secure" mapstructure:"is_secure"`
	TxData   RippleTxData `toml:"transaction" mapstructure:"transaction"`
}

// RippleTxData is used for api command to send coin
type RippleTxData struct {
	Account string `toml:"sender_account" mapstructure:"sender_account"`
	Secret  string `toml:"sender_secret" mapstructure:"sender_secret"`
}

// Logger logger info
type Logger struct {
	Service  string `toml:"service" mapstructure:"service" validate:"required"`
	Env      string `toml:"env" mapstructure:"env" validate:"oneof=dev prod custom"`
	Level    string `toml:"level" mapstructure:"level" validate:"required"`
	IsLogger bool   `toml:"is_logger" mapstructure:"is_logger"`
}

// Tracer is open tracing
type Tracer struct {
	Type    string       `toml:"type" mapstructure:"type" validate:"oneof=none jaeger datadog"`
	Jaeger  TracerDetail `toml:"jaeger" mapstructure:"jaeger"`
	Datadog TracerDetail `toml:"datadog" mapstructure:"datadog"`
}

// TracerDetail includes specific service config
type TracerDetail struct {
	ServiceName         string  `toml:"service_name" mapstructure:"service_name"`
	CollectorEndpoint   string  `toml:"collector_endpoint" mapstructure:"collector_endpoint"`
	SamplingProbability float64 `toml:"sampling_probability" mapstructure:"sampling_probability"`
	IsDebug             bool    `toml:"is_debug" mapstructure:"is_debug"`
}

// MySQL info
type MySQL struct {
	Host  string `toml:"host" mapstructure:"host" validate:"required"`
	DB    string `toml:"dbname" mapstructure:"dbname" validate:"required"`
	User  string `toml:"user" mapstructure:"user" validate:"required"`
	Pass  string `toml:"pass" mapstructure:"pass" validate:"required"`
	Debug bool   `toml:"debug" mapstructure:"debug"`
}

// FilePath if file path group
type FilePath struct {
	Tx         string `toml:"tx" mapstructure:"tx" validate:"required"`
	Address    string `toml:"address" mapstructure:"address" validate:"required"`
	FullPubKey string `toml:"full_pubkey" mapstructure:"full_pubkey" validate:"required"`
}

// PubKeyFile saved pubKey file path which is used when import/export file
type PubKeyFile struct {
	BasePath string `toml:"base_path" mapstructure:"base_path" validate:"required"`
}

// AddressFile saved address file path which is used when import/export file
type AddressFile struct {
	BasePath string `toml:"base_path" mapstructure:"base_path" validate:"required"`
}
