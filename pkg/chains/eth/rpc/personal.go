package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// ImportRawKey imports the given unencrypted private key (hex string) into the key store,
// encrypting it with the passphrase.
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_importrawkey
func (c *rpcClient) ImportRawKey(ctx context.Context, hexKey, passPhrase string) (string, error) {
	if strings.HasPrefix(hexKey, "0x") {
		hexKey = strings.TrimLeft(hexKey, "0x")
	}

	var address string
	err := c.caller.CallContext(ctx, &address, "personal_importRawKey", hexKey, passPhrase)
	if err != nil {
		if strings.Contains(err.Error(), "account already exists") {
			logger.Warn(err.Error())
		} else {
			return "", fmt.Errorf("fail to call client.CallContext(personal_importRawKey): %w", err)
		}
	}

	return address, nil
}

// ListAccounts returns all Ethereum account addresses stored in the key store.
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_listaccounts
func (c *rpcClient) ListAccounts(ctx context.Context) ([]string, error) {
	var accounts []string
	err := c.caller.CallContext(ctx, &accounts, "personal_listAccounts")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(personal_listAccounts): %w", err)
	}

	return accounts, nil
}

// NewAccount generates a new private key and stores it in the key store directory.
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_newaccount
func (c *rpcClient) NewAccount(ctx context.Context, passphrase string) (string, error) {
	var address string
	err := c.caller.CallContext(ctx, &address, "personal_newAccount", passphrase)
	if err != nil {
		return "", fmt.Errorf("fail to call rpc.CallContext(personal_newAccount): %w", err)
	}
	return address, nil
}

// LockAccount removes the private key with the given address from memory.
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_lockaccount
func (c *rpcClient) LockAccount(ctx context.Context, hexAddr string) error {
	err := c.caller.CallContext(ctx, nil, "personal_lockAccount", hexAddr)
	if err != nil {
		return fmt.Errorf("fail to call rpc.CallContext(personal_lockAccount): %w", err)
	}
	return err
}

// UnlockAccount decrypts the key with the given address for the specified duration (seconds).
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_unlockaccount
func (c *rpcClient) UnlockAccount(ctx context.Context, hexAddr, passphrase string, duration uint64) (bool, error) {
	var ret bool
	err := c.caller.CallContext(ctx, &ret, "personal_unlockAccount", hexAddr, passphrase, duration)
	if err != nil {
		return false, fmt.Errorf("fail to call rpc.CallContext(personal_unlockAccount): %w", err)
	}
	return ret, nil
}
