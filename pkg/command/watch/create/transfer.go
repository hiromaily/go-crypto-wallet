package create

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runTransfer(wallet wallets.Watcher, account1, account2 string, amount, fee float64) error {
	// validator
	if !account.ValidateAccountType(account1) {
		return errors.New("account option [-account1] is invalid")
	}
	if !account.ValidateAccountType(account2) {
		return errors.New("account option [-account2] is invalid")
	}
	// This logic should be implemented in wallet.CreateTransferTx()
	// if !account.NotAllow(account1, []account.AccountType{
	//	account.AccountTypeAuthorization, account.AccountTypeClient}) {
	//	return fmt.Errorf(
	//		"account1: %s/%s is not allowed",
	//		account.AccountTypeAuthorization, account.AccountTypeClient)
	//}
	// if !account.NotAllow(account2, []account.AccountType{
	//	account.AccountTypeAuthorization, account.AccountTypeClient}) {
	//	return fmt.Errorf(
	//		"account2: %s/%s is not allowed",
	//		account.AccountTypeAuthorization, account.AccountTypeClient)
	//}
	// if amount == 0{
	//	return fmt.Errorf("amount option [-amount] is invalid")
	//}

	hex, fileName, err := wallet.CreateTransferTx(
		account.AccountType(account1),
		account.AccountType(account2),
		amount,
		fee)
	if err != nil {
		return fmt.Errorf("fail to call CreateTransferTx() %w", err)
	}
	if (wallet.CoinTypeCode() != coin.ETH && wallet.CoinTypeCode() != coin.XRP) && hex == "" {
		fmt.Println("No utxo")
		return nil
	}
	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[fileName]: %s\n", hex, fileName)

	return nil
}
