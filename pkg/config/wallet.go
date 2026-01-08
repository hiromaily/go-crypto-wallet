package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

// WalletRoot wallet root config
type WalletRoot struct {
	//nolint:lll,revive
	KeyType domainKey.KeyType `toml:"key_type" yaml:"key_type" mapstructure:"key_type" validate:"omitempty,oneof=bip44 bip49 bip84 bip86 musig2"`
	//nolint:lll,revive
	AddressType  domainAddress.AddrType  `toml:"address_type" yaml:"address_type" mapstructure:"address_type" validate:"oneof=legacy p2sh-segwit bech32 bch-cashaddr taproot"`
	CoinTypeCode domainCoin.CoinTypeCode `toml:"coin_type" yaml:"coin_type" mapstructure:"coin_type"`
	Bitcoin      Bitcoin                 `toml:"bitcoin" yaml:"bitcoin" mapstructure:"bitcoin"`
	Ethereum     Ethereum                `toml:"ethereum" yaml:"ethereum" mapstructure:"ethereum"`
	Ripple       Ripple                  `toml:"ripple" yaml:"ripple" mapstructure:"ripple"`
	Logger       Logger                  `toml:"logger" yaml:"logger" mapstructure:"logger"`
	Tracer       Tracer                  `toml:"tracer" yaml:"tracer" mapstructure:"tracer"`
	MySQL        MySQL                   `toml:"mysql" yaml:"mysql" mapstructure:"mysql"`
	FilePath     FilePath                `toml:"file_path" yaml:"file_path" mapstructure:"file_path"`
}

// Bitcoin information
type Bitcoin struct {
	Host       string `toml:"host" yaml:"host" mapstructure:"host" validate:"required"`
	User       string `toml:"user" yaml:"user" mapstructure:"user" validate:"required"`
	Pass       string `toml:"pass" yaml:"pass" mapstructure:"pass" validate:"required"`
	PostMode   bool   `toml:"http_post_mode" yaml:"http_post_mode" mapstructure:"http_post_mode"`
	DisableTLS bool   `toml:"disable_tls" yaml:"disable_tls" mapstructure:"disable_tls"`
	//nolint:lll,revive
	NetworkType string `toml:"network_type" yaml:"network_type" mapstructure:"network_type" validate:"oneof=mainnet testnet3 regtest signet"`

	Block BitcoinBlock `toml:"block" yaml:"block" mapstructure:"block"`
	Fee   BitcoinFee   `toml:"fee" yaml:"fee" mapstructure:"fee"`
}

// BitcoinBlock block information of Bitcoin.
//
// Note: This field is only required for watch-only wallets.
// Keygen and signature wallets do not require this value, so validation
// is conditionally applied based on wallet type.
type BitcoinBlock struct {
	ConfirmationNum uint64 `toml:"confirmation_num" yaml:"confirmation_num" mapstructure:"confirmation_num"`
}

// BitcoinFee range of adjustment calculated fee when sending coin
type BitcoinFee struct {
	AdjustmentMin float64 `toml:"adjustment_min" yaml:"adjustment_min" mapstructure:"adjustment_min"`
	AdjustmentMax float64 `toml:"adjustment_max" yaml:"adjustment_max" mapstructure:"adjustment_max"`
}

// Ethereum information
type Ethereum struct {
	Host       string `toml:"host" yaml:"host" mapstructure:"host" validate:"required"`
	IPCPath    string `toml:"ipc_path" yaml:"ipc_path" mapstructure:"ipc_path"`
	Port       int    `toml:"port" yaml:"port" mapstructure:"port" validate:"required"`
	DisableTLS bool   `toml:"disable_tls" yaml:"disable_tls" mapstructure:"disable_tls"`
	//nolint:lll,revive
	NetworkType string `toml:"network_type" yaml:"network_type" mapstructure:"network_type" validate:"oneof=mainnet goerli rinkeby ropsten"`
	KeyDirName  string `toml:"keydir" yaml:"keydir" mapstructure:"keydir"`
	//nolint:lll,revive
	ConfirmationNum uint64                          `toml:"confirmation_num" yaml:"confirmation_num" mapstructure:"confirmation_num"`
	ERC20Token      domainCoin.ERC20Token           `toml:"erc20_token" yaml:"erc20_token" mapstructure:"erc20_token"`
	ERC20s          map[domainCoin.ERC20Token]ERC20 `toml:"erc20s" yaml:"erc20s" mapstructure:"erc20s"`
}

// ERC20 information
type ERC20 struct {
	Symbol          string `toml:"symbol" yaml:"symbol" mapstructure:"symbol"`
	Name            string `toml:"name" yaml:"name" mapstructure:"name"`
	ContractAddress string `toml:"contract_address" yaml:"contract_address" mapstructure:"contract_address"`
	MasterAddress   string `toml:"master_address" yaml:"master_address" mapstructure:"master_address"`
	Decimals        int    `toml:"decimals" yaml:"decimals" mapstructure:"decimals"`
}

// Ripple information
type Ripple struct {
	//nolint:lll,revive
	WebsocketPublicURL string `toml:"websocket_public_url" yaml:"websocket_public_url" mapstructure:"websocket_public_url"`
	WebsocketAdminURL  string `toml:"websocket_admin_url" yaml:"websocket_admin_url" mapstructure:"websocket_admin_url"`
	//nolint:lll,revive
	NetworkType string    `toml:"network_type" yaml:"network_type" mapstructure:"network_type" validate:"oneof=mainnet testnet devnet"`
	API         RippleAPI `toml:"api" yaml:"api" mapstructure:"api"`
}

// RippleAPI is ripple-lib server info
type RippleAPI struct {
	URL      string       `toml:"url" yaml:"url" mapstructure:"url"`
	IsSecure bool         `toml:"is_secure" yaml:"is_secure" mapstructure:"is_secure"`
	TxData   RippleTxData `toml:"transaction" yaml:"transaction" mapstructure:"transaction"`
}

// RippleTxData is used for api command to send coin
type RippleTxData struct {
	Account string `toml:"sender_account" yaml:"sender_account" mapstructure:"sender_account"`
	Secret  string `toml:"sender_secret" yaml:"sender_secret" mapstructure:"sender_secret"`
}

// Logger logger info
type Logger struct {
	Service  string `toml:"service" yaml:"service" mapstructure:"service" validate:"required"`
	Format   string `toml:"format" yaml:"format" mapstructure:"format" validate:"oneof=json console"`
	Level    string `toml:"level" yaml:"level" mapstructure:"level" validate:"required"`
	IsLogger bool   `toml:"is_logger" yaml:"is_logger" mapstructure:"is_logger"`
}

// Tracer is open tracing
type Tracer struct {
	Type    string       `toml:"type" yaml:"type" mapstructure:"type" validate:"oneof=none jaeger datadog"`
	Jaeger  TracerDetail `toml:"jaeger" yaml:"jaeger" mapstructure:"jaeger"`
	Datadog TracerDetail `toml:"datadog" yaml:"datadog" mapstructure:"datadog"`
}

// TracerDetail includes specific service config
type TracerDetail struct {
	ServiceName       string `toml:"service_name" yaml:"service_name" mapstructure:"service_name"`
	CollectorEndpoint string `toml:"collector_endpoint" yaml:"collector_endpoint" mapstructure:"collector_endpoint"`
	//nolint:lll,revive
	SamplingProbability float64 `toml:"sampling_probability" yaml:"sampling_probability" mapstructure:"sampling_probability"`
	IsDebug             bool    `toml:"is_debug" yaml:"is_debug" mapstructure:"is_debug"`
}

// MySQL info
type MySQL struct {
	Host  string `toml:"host" yaml:"host" mapstructure:"host" validate:"required"`
	DB    string `toml:"dbname" yaml:"dbname" mapstructure:"dbname" validate:"required"`
	User  string `toml:"user" yaml:"user" mapstructure:"user" validate:"required"`
	Pass  string `toml:"pass" yaml:"pass" mapstructure:"pass" validate:"required"`
	Debug bool   `toml:"debug" yaml:"debug" mapstructure:"debug"`
}

// FilePath if file path group
type FilePath struct {
	Tx         string `toml:"tx" yaml:"tx" mapstructure:"tx" validate:"required"`
	Address    string `toml:"address" yaml:"address" mapstructure:"address" validate:"required"`
	FullPubKey string `toml:"full_pubkey" yaml:"full_pubkey" mapstructure:"full_pubkey" validate:"required"`
}

// PubKeyFile saved pubKey file path which is used when import/export file
type PubKeyFile struct {
	BasePath string `toml:"base_path" yaml:"base_path" mapstructure:"base_path" validate:"required"`
}

// AddressFile saved address file path which is used when import/export file
type AddressFile struct {
	BasePath string `toml:"base_path" yaml:"base_path" mapstructure:"base_path" validate:"required"`
}

// NewWallet creates wallet config by loading it from the specified file.
//
// This function:
//  1. Validates the file path
//  2. Loads the configuration file (supports both TOML and YAML formats)
//  3. Validates the configuration structure based on wallet type and coin type
//  4. Returns the validated WalletRoot configuration
//
// The function automatically detects the configuration format based on file extension:
//   - .toml -> TOML format
//   - .yaml, .yml -> YAML format
//
// Returns an error if:
//   - File path is empty
//   - File extension is not supported
//   - File cannot be read
//   - Unmarshaling fails
//   - Validation fails
func NewWallet(file string, wtype domainWallet.WalletType, coinTypeCode domainCoin.CoinTypeCode) (*WalletRoot, error) {
	if file == "" {
		return nil, errors.New("wallet config file path cannot be empty")
	}

	var conf WalletRoot
	if err := loadConfig(file, &conf); err != nil {
		return nil, err
	}

	if err := conf.validate(wtype, coinTypeCode); err != nil {
		return nil, fmt.Errorf("wallet config validation failed: %w", err)
	}

	return &conf, nil
}

// loadWallet loads wallet configuration from a TOML file.
// This function is kept for backward compatibility with tests.
func loadWallet(path string) (*WalletRoot, error) {
	return loadTOML[WalletRoot](path)
}

// validate validates the wallet configuration structure based on wallet type and coin type.
func (c *WalletRoot) validate(wtype domainWallet.WalletType, coinTypeCode domainCoin.CoinTypeCode) error {
	validate := validator.New()

	switch coinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		if err := validate.StructExcept(c, "Ethereum", "Ripple"); err != nil {
			return err
		}
		switch wtype {
		case domainWallet.WalletTypeWatchOnly:
			if c.Bitcoin.Block.ConfirmationNum == 0 {
				return errors.New("block ConfirmationNum is required in toml file")
			}
		case domainWallet.WalletTypeKeyGen, domainWallet.WalletTypeSign:
			// No additional validation needed
		default:
		}
	case domainCoin.ETH, domainCoin.ERC20:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ripple"); err != nil {
			return err
		}
	case domainCoin.XRP:
		if err := validate.StructExcept(c, "AddressType", "Bitcoin", "Ethereum"); err != nil {
			return err
		}
	case domainCoin.LTC, domainCoin.HYT:
		// Not implemented yet
	default:
	}

	return nil
}

// ValidateERC20 validates that the specified ERC20 token is configured.
func (c *WalletRoot) ValidateERC20(token domainCoin.ERC20Token) error {
	if _, ok := c.Ethereum.ERC20s[token]; !ok {
		return fmt.Errorf("erc20 token information for [%s] is required", token.String())
	}
	return nil
}
