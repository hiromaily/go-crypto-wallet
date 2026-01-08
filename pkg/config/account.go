package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// AccountRoot represents account configuration loaded from TOML files.
//
// This structure is placed in pkg/config (not internal/infrastructure/config) for the following reasons:
//
// 1. Used by both pkg/di and internal/di:
//   - pkg/di uses AccountRoot to initialize reusable components
//   - internal/di uses AccountRoot to initialize application-specific components
//   - Both DI containers need access to this configuration during initialization
//
// 2. Architectural compliance:
//   - pkg/ packages MUST NOT import from internal/ (pkg/AGENTS.md)
//   - If AccountRoot were in internal/infrastructure/config, pkg/di would violate this rule
//   - Moving to pkg/config ensures pkg/di can use it without violating architectural principles
//
// 3. Initialization phase usage:
//   - Configuration is loaded in cmd/ (main entry points) before DI container creation
//   - Both pkg/di and internal/di receive AccountRoot as a parameter during initialization
//   - This is similar to WalletRoot, which is also in pkg/config
//
// 4. Domain types usage:
//   - AccountRoot uses domain types (AccountType, AuthType) which is correct
//   - Dependency direction: pkg/config → domain (allowed)
//   - This represents configuration data that maps to domain concepts
//
// Note: While AccountRoot is application-specific, it's placed in pkg/ because it's
// shared between pkg/di and internal/di during the initialization phase.
type AccountRoot struct {
	Types []domainAccount.AccountType `toml:"types" yaml:"types" mapstructure:"types"`
	//nolint:revive
	DepositReceiver domainAccount.AccountType `toml:"deposit_receiver" yaml:"deposit_receiver" mapstructure:"deposit_receiver"`

	PaymentSender domainAccount.AccountType `toml:"payment_sender" yaml:"payment_sender" mapstructure:"payment_sender"`
	Multisigs     []AccountMultisig         `toml:"multisig" yaml:"multisig" mapstructure:"multisig"`
}

// AccountMultisig represents multisig account configuration.
//
// This structure is placed in pkg/config for the same reasons as AccountRoot:
// - Used by both pkg/di and internal/di during initialization
// - Must be accessible from pkg/ without importing internal/
// - Represents configuration data that maps to domain concepts
//
// The structure contains TOML and YAML tags for file deserialization, which is an
// infrastructure concern, but the structure itself is shared configuration
// that needs to be accessible from both pkg/ and internal/ layers.
type AccountMultisig struct {
	Type      domainAccount.AccountType `toml:"type" yaml:"type" mapstructure:"type"`
	Required  int                       `toml:"required" yaml:"required" mapstructure:"required"`
	AuthUsers []domainAccount.AuthType  `toml:"auth_users" yaml:"auth_users" mapstructure:"auth_users"`
}

// NewAccount loads and validates account configuration from a file.
//
// This function is placed in pkg/config (not internal/infrastructure/config) because:
//
// 1. Used during initialization phase:
//   - Called from cmd/ (main entry points) before DI container creation
//   - Both pkg/di and internal/di need AccountRoot during initialization
//   - Must be accessible without importing internal/ packages
//
// 2. Configuration loading is a shared concern:
//   - Similar to NewWallet() which is also in pkg/config
//   - Both functions handle configuration file loading and validation
//   - This is infrastructure-level file I/O, but needed at initialization
//
// 3. Architectural compliance:
//   - pkg/ packages cannot import from internal/
//   - If this were in internal/infrastructure/config, cmd/ would need to import internal/
//   - While cmd/ can import internal/, keeping config loading in pkg/ is cleaner
//
// The function performs:
//   - File path validation
//   - Configuration file loading (supports both TOML and YAML formats)
//   - Structure validation using go-playground/validator
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
func NewAccount(file string) (*AccountRoot, error) {
	if file == "" {
		return nil, errors.New("account config file path cannot be empty")
	}

	var conf AccountRoot
	if err := loadConfig(file, &conf); err != nil {
		return nil, err
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("account config validation failed: %w", err)
	}

	return &conf, nil
}

// validate validates the account configuration structure using go-playground/validator.
func (c *AccountRoot) validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

// NewMultisigConfig converts AccountMultisig config to domain MultisigConfig.
//
// This function is placed in pkg/config (not internal/infrastructure/config/account) because:
//
// 1. Co-location with configuration structures:
//   - AccountMultisig and AccountRoot are defined in pkg/config
//   - Having the conversion function in the same package improves code discoverability
//   - All account configuration-related code is in one place
//
// 2. Used during initialization phase:
//   - Called from internal/di during container initialization
//   - Both pkg/di and internal/di may need this function
//   - Placing it in pkg/config makes it accessible from both layers
//
// 3. Dependency management:
//   - pkg/config already depends on internal/domain/account (AccountMultisig uses domain types)
//   - This function also depends on internal/domain/account (calls domainAccount.NewMultisigConfig)
//   - No circular dependencies are introduced
//
// 4. Separation of concerns:
//   - pkg/config: Configuration structures and conversion to domain entities
//   - Domain layer: Domain entity factory (NewMultisigConfig) using domain types
//   - This function bridges configuration (pkg/config) and domain entities
//
// The function:
//  1. Takes []AccountMultisig (from pkg/config)
//  2. Converts to domain types (AccountType, AuthType)
//  3. Calls domainAccount.NewMultisigConfig() to create the domain entity
//
// This pattern follows Clean Architecture: Configuration layer converts config structures
// to domain types, then uses domain factories to create domain entities.
func NewMultisigConfig(confMultisig []AccountMultisig) *domainAccount.MultisigConfig {
	if confMultisig == nil {
		return domainAccount.NewMultisigConfig(nil)
	}

	accountMap := make(map[domainAccount.AccountType]map[int][]domainAccount.AuthType, len(confMultisig))
	for _, val := range confMultisig {
		accountMap[val.Type] = map[int][]domainAccount.AuthType{
			val.Required: val.AuthUsers,
		}
	}
	return domainAccount.NewMultisigConfig(accountMap)
}
