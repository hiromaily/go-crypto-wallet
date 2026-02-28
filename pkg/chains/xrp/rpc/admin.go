package rpc

import (
	"context"
	"fmt"
)

// Note: Admin commands require a connection to rippled on an admin host/port.
// https://xrpl.org/key-generation-methods.html

// ValidationCreateRequest is the request payload for the validation_create command.
type ValidationCreateRequest struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
	Secret  string `json:"secret"`
}

// ResponseValidationCreate is the wire-format response of the validation_create command.
type ResponseValidationCreate struct {
	Result struct {
		Status              string `json:"status"`
		ValidationKey       string `json:"validation_key"`
		ValidationPublicKey string `json:"validation_public_key"`
		ValidationSeed      string `json:"validation_seed"`
	} `json:"result"`
	Error string `json:"error,omitempty"`
}

// WalletProposeWithKeyRequest is the request payload for wallet_propose with an explicit key type.
type WalletProposeWithKeyRequest struct {
	Command string `json:"command"`
	Seed    string `json:"seed"`
	KeyType string `json:"key_type"`
}

// WalletProposeRequest is the request payload for wallet_propose using a passphrase.
type WalletProposeRequest struct {
	Command    string `json:"command"`
	Passphrase string `json:"passphrase"`
}

// ResponseWalletPropose is the wire-format response of the wallet_propose command.
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

// ValidationCreate calls the validation_create admin WebSocket command and returns the raw wire response.
func ValidationCreate(ctx context.Context, caller WSCaller, secret string) (*ResponseValidationCreate, error) {
	req := &ValidationCreateRequest{
		ID:      0,
		Command: "validation_create",
		Secret:  secret,
	}
	var res ResponseValidationCreate
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsAdmin.Call(validation_create): %w", err)
	}
	return &res, nil
}

// WalletProposeWithKey calls wallet_propose with a seed and explicit key type.
func WalletProposeWithKey(
	ctx context.Context, caller WSCaller, seed string, keyType KeyType,
) (*ResponseWalletPropose, error) {
	req := &WalletProposeWithKeyRequest{
		Command: "wallet_propose",
		Seed:    seed,
		KeyType: keyType.String(),
	}
	var res ResponseWalletPropose
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsAdmin.Call(wallet_propose): %w", err)
	}
	return &res, nil
}

// WalletPropose calls wallet_propose with a passphrase.
// The result is deterministic: the same passphrase always produces the same key pair.
func WalletPropose(ctx context.Context, caller WSCaller, passphrase string) (*ResponseWalletPropose, error) {
	req := &WalletProposeRequest{
		Command:    "wallet_propose",
		Passphrase: passphrase,
	}
	var res ResponseWalletPropose
	if err := caller.Call(ctx, req, &res); err != nil {
		return nil, fmt.Errorf("fail to call wsAdmin.Call(wallet_propose): %w", err)
	}
	return &res, nil
}
