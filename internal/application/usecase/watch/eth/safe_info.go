package eth

import (
	"context"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
)

type ethSafeInfoUseCase struct {
	safeInfoReader apieth.SafeInfoReader
}

// NewETHSafeInfoUseCase constructs the use case for retrieving on-chain Safe contract state.
func NewETHSafeInfoUseCase(safeInfoReader apieth.SafeInfoReader) watchusecase.ETHSafeInfoUseCase {
	return &ethSafeInfoUseCase{safeInfoReader: safeInfoReader}
}

// Execute fetches Safe contract state and maps it to the output DTO.
func (u *ethSafeInfoUseCase) Execute(
	ctx context.Context,
	input watchusecase.ETHSafeInfoInput,
) (watchusecase.ETHSafeInfoOutput, error) {
	info, err := u.safeInfoReader.GetSafeInfo(ctx, input.SafeAddress)
	if err != nil {
		return watchusecase.ETHSafeInfoOutput{},
			fmt.Errorf("failed to get Safe info: %w", err)
	}

	return watchusecase.ETHSafeInfoOutput{
		Owners:     info.Owners,
		Threshold:  info.Threshold,
		Nonce:      info.Nonce,
		BalanceWei: info.Balance.String(),
	}, nil
}
