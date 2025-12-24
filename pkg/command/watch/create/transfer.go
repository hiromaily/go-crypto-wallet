package create

import (
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runTransfer(wallet wallets.Watcher, account1, account2 string, amount, fee float64) error {
	// validator
	if !domainAccount.ValidateAccountType(account1) {
		return errors.New("account option [-account1] is invalid")
	}
	if !domainAccount.ValidateAccountType(account2) {
		return errors.New("account option [-account2] is invalid")
	}
	// This logic should be implemented in wallet.CreateTransferTx()
	// if !domainAccount.NotAllow(account1, []domainAccount.AccountType{
	//	domainAccount.AccountTypeAuthorization, domainAccount.AccountTypeClient}) {
	//	return fmt.Errorf(
	//		"account1: %s/%s is not allowed",
	//		domainAccount.AccountTypeAuthorization, domainAccount.AccountTypeClient)
	//}
	// if !domainAccount.NotAllow(account2, []domainAccount.AccountType{
	//	domainAccount.AccountTypeAuthorization, domainAccount.AccountTypeClient}) {
	//	return fmt.Errorf(
	//		"account2: %s/%s is not allowed",
	//		domainAccount.AccountTypeAuthorization, domainAccount.AccountTypeClient)
	//}
	// if amount == 0{
	//	return fmt.Errorf("amount option [-amount] is invalid")
	//}

	hex, fileName, err := wallet.CreateTransferTx(
		domainAccount.AccountType(account1),
		domainAccount.AccountType(account2),
		amount,
		fee)
	if err != nil {
		return fmt.Errorf("fail to call CreateTransferTx() %w", err)
	}
	if (wallet.CoinTypeCode() != domainCoin.ETH && wallet.CoinTypeCode() != domainCoin.XRP) && hex == "" {
		fmt.Println("No utxo")
		return nil
	}
	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[fileName]: %s\n", hex, fileName)

	return nil
}
