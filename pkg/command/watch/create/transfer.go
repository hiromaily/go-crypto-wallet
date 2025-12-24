package create

import (
	"context"
	"errors"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
)

func runTransfer(container di.Container, account1, account2 string, amount, fee float64) error {
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

	// Get use case from container
	useCase := container.NewWatchCreateTransactionUseCase().(watchusecase.CreateTransactionUseCase)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountType(account1),
		ReceiverAccount: domainAccount.AccountType(account2),
		Amount:          amount,
		AdjustmentFee:   fee,
	})
	if err != nil {
		return fmt.Errorf("fail to create transfer transaction: %w", err)
	}

	if output.TransactionHex == "" {
		fmt.Println("No utxo")
		return nil
	}

	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[fileName]: %s\n", output.TransactionHex, output.FileName)

	return nil
}
