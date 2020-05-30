package rippleapi

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	pb "github.com/hiromaily/ripple-lib-proto/pb/go/rippleapi"
)

// PrepareTransaction calls PrepareTransaction API
func (r *RippleAPI) PrepareTransaction(senderAccount, receiverAccount string, amount float64) error {
	ctx := context.Background()
	req := &pb.RequestPrepareTransaction{
		TxType:          pb.TX_PAYMENT,
		SenderAccount:   senderAccount,
		Amount:          amount,
		ReceiverAccount: receiverAccount,
		Instructions:    nil,
	}

	res, err := r.client.PrepareTransaction(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, "fail to call client.PrepareTransaction()")
	}
	r.logger.Debug("response", zap.Any("response", res))

	return nil
}
