package eth

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// ImportRawKey Imports the given unencrypted private key (hex string) into the key store,
// encrypting it with the passphrase
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_importrawkey
func (e *Ethereum) ImportRawKey(ctx context.Context, hexKey, passPhrase string) (string, error) {
	return e.pkgrpc.ImportRawKey(ctx, hexKey, passPhrase)
}

// ListAccounts returns all the Ethereum account addresses of all keys in the key store
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_listaccounts
func (e *Ethereum) ListAccounts(ctx context.Context) ([]string, error) {
	return e.pkgrpc.ListAccounts(ctx)
}

// NewAccount generates a new private key and stores it in the key store directory
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_newaccount
func (e *Ethereum) NewAccount(
	ctx context.Context, passphrase string, _ domainAccount.AccountType,
) (string, error) {
	return e.pkgrpc.NewAccount(ctx, passphrase)
}

// LockAccount removes the private key with given address from memory
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_lockaccount
func (e *Ethereum) LockAccount(ctx context.Context, hexAddr string) error {
	return e.pkgrpc.LockAccount(ctx, hexAddr)
}

// UnlockAccount decrypts the key with the given address from the key store.
//
//	duration: second
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_unlockaccount
func (e *Ethereum) UnlockAccount(ctx context.Context, hexAddr, passphrase string, duration uint64) (bool, error) {
	return e.pkgrpc.UnlockAccount(ctx, hexAddr, passphrase, duration)
}
