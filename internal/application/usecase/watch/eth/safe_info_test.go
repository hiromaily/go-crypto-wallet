package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
)

const siTestSafeAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

type ethSafeInfoDeps struct {
	safeInfoReader *ethapiamocks.MockSafeInfoReader
}

func newEthSafeInfoDeps(t *testing.T) *ethSafeInfoDeps {
	t.Helper()
	return &ethSafeInfoDeps{
		safeInfoReader: ethapiamocks.NewMockSafeInfoReader(t),
	}
}

func newEthSafeInfoUseCase(deps *ethSafeInfoDeps) watchusecase.ETHSafeInfoUseCase {
	return watchusecaseeth.NewETHSafeInfoUseCase(deps.safeInfoReader)
}

func TestETHSafeInfoUseCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   watchusecase.ETHSafeInfoInput
		setup   func(deps *ethSafeInfoDeps)
		wantErr bool
		checkFn func(t *testing.T, out watchusecase.ETHSafeInfoOutput)
	}{
		{
			name:  "returns Safe info successfully",
			input: watchusecase.ETHSafeInfoInput{SafeAddress: siTestSafeAddr},
			setup: func(deps *ethSafeInfoDeps) {
				info := &apieth.SafeInfo{
					Owners: []string{
						"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
						"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
					},
					Threshold: 2,
					Nonce:     7,
					Balance:   big.NewInt(500000000000000000), // 0.5 ETH
				}
				deps.safeInfoReader.EXPECT().GetSafeInfo(
					context.Background(), siTestSafeAddr,
				).Return(info, nil).Once()
			},
			wantErr: false,
			checkFn: func(t *testing.T, out watchusecase.ETHSafeInfoOutput) {
				t.Helper()
				require.Len(t, out.Owners, 2)
				assert.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", out.Owners[0])
				assert.Equal(t, "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", out.Owners[1])
				assert.Equal(t, uint64(2), out.Threshold)
				assert.Equal(t, uint64(7), out.Nonce)
				assert.Equal(t, "500000000000000000", out.BalanceWei)
			},
		},
		{
			name:  "Safe with zero balance",
			input: watchusecase.ETHSafeInfoInput{SafeAddress: siTestSafeAddr},
			setup: func(deps *ethSafeInfoDeps) {
				info := &apieth.SafeInfo{
					Owners:    []string{"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
					Threshold: 1,
					Nonce:     0,
					Balance:   big.NewInt(0),
				}
				deps.safeInfoReader.EXPECT().GetSafeInfo(
					context.Background(), siTestSafeAddr,
				).Return(info, nil).Once()
			},
			wantErr: false,
			checkFn: func(t *testing.T, out watchusecase.ETHSafeInfoOutput) {
				t.Helper()
				assert.Equal(t, uint64(1), out.Threshold)
				assert.Equal(t, uint64(0), out.Nonce)
				assert.Equal(t, "0", out.BalanceWei)
			},
		},
		{
			name:  "error on GetSafeInfo failure",
			input: watchusecase.ETHSafeInfoInput{SafeAddress: siTestSafeAddr},
			setup: func(deps *ethSafeInfoDeps) {
				deps.safeInfoReader.EXPECT().GetSafeInfo(
					context.Background(), siTestSafeAddr,
				).Return(nil, errors.New("RPC error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newEthSafeInfoDeps(t)
			if tt.setup != nil {
				tt.setup(deps)
			}

			uc := newEthSafeInfoUseCase(deps)
			out, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.checkFn != nil {
				tt.checkFn(t, out)
			}
		})
	}
}
