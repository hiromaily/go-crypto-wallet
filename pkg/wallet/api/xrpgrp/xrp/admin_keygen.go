package xrp

import (
	"context"

	"github.com/pkg/errors"
)

// Note: Admin commands are available only if you connect to rippled on a host and port that the rippled.cfg file identifies as admin

// https://xrpl.org/key-generation-methods.html

// Assign a Regular Key Pair
// https://xrpl.org/assign-a-regular-key-pair.html
// https://github.com/ripple/ripple-keypairs

// ValidationCreate is request data for validation_create method
type ValidationCreate struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
	Secret  string `json:"secret"`
}

// ResponseValidationCreate is response data for validation_create method
type ResponseValidationCreate struct {
	Result struct {
		Status              string `json:"status"`
		ValidationKey       string `json:"validation_key"`
		ValidationPublicKey string `json:"validation_public_key"`
		ValidationSeed      string `json:"validation_seed"`
	} `json:"result"`
	Error string `json:"error,omitempty"`
}

// ValidationCreate calls validation_create method
func (r *Ripple) ValidationCreate(secret string) (*ResponseValidationCreate, error) {
	if r.wsAdmin == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	req := ValidationCreate{
		ID:      0,
		Command: "validation_create",
		Secret:  secret,
	}
	var res ResponseValidationCreate
	if err := r.wsAdmin.Call(context.Background(), &req, &res); err != nil {
		return nil, errors.Wrap(err, "fail to call wsAdmin.Call(validation_create)")
	}
	return &res, nil
}

// WalletProposeWithKey is request data for wallet_propose method
type WalletProposeWithKey struct {
	Command string `json:"command"`
	Seed    string `json:"seed"`
	KeyType string `json:"key_type"`
}

// WalletPropose is request data for wallet_propose method
type WalletPropose struct {
	Command    string `json:"command"`
	Passphrase string `json:"passphrase"`
}

// ResponseWalletPropose is response data for wallet_propose method
type ResponseWalletPropose struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Result struct {
		AccountID     string `json:"account_id"`
		KeyType       string `json:"key_type"`
		MasterKey     string `json:"master_key"` // DEPRECATED
		MasterSeed    string `json:"master_seed"`
		MasterSeedHex string `json:"master_seed_hex"`
		PublicKey     string `json:"public_key"`
		PublicKeyHex  string `json:"public_key_hex"`
	} `json:"result"`
	Error string `json:"error,omitempty"`
}

// WalletProposeWithKey calls wallet_propose method
func (r *Ripple) WalletProposeWithKey(seed string, keyType XRPKeyType) (*ResponseWalletPropose, error) {
	if r.wsAdmin == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	req := WalletProposeWithKey{
		Command: "wallet_propose",
		Seed:    seed,
		KeyType: keyType.String(),
	}
	var res ResponseWalletPropose
	if err := r.wsAdmin.Call(context.Background(), &req, &res); err != nil {
		return nil, errors.Wrap(err, "fail to call wsAdmin.Call(wallet_propose)")
	}
	return &res, nil
}

// WalletPropose calls wallet_propose method
// - result is same as long as using same passphrase
func (r *Ripple) WalletPropose(passphrase string) (*ResponseWalletPropose, error) {
	if r.wsAdmin == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	req := WalletPropose{
		Command:    "wallet_propose",
		Passphrase: passphrase,
	}
	var res ResponseWalletPropose
	if err := r.wsAdmin.Call(context.Background(), &req, &res); err != nil {
		return nil, errors.Wrap(err, "fail to call wsAdmin.Call(wallet_propose)")
	}
	return &res, nil
}
