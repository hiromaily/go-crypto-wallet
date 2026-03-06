package public

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // ─── AccountChannels ──────────────────────────────────────────────────────────

// func TestAccountChannels_HappyPath(t *testing.T) {
// 	t.Parallel()
// 	caller := &mockWSCaller{
// 		callFn: func(_ context.Context, req, res any) error {
// 			r := req.(*AccountChannelsRequest)
// 			assert.Equal(t, "account_channels", r.Command)
// 			assert.Equal(t, "rSender", r.Account)
// 			assert.Equal(t, "rReceiver", r.DestinationAccount)
// 			p := res.(*ResponseAccountChannels)
// 			p.Status = "success"
// 			p.Result.Account = "rSender"
// 			return nil
// 		},
// 	}
// 	res, err := AccountChannels(context.Background(), caller, "rSender", "rReceiver")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// 	assert.Equal(t, "success", res.Status)
// 	assert.Equal(t, "rSender", res.Result.Account)
// }

// func TestAccountChannels_Error(t *testing.T) {
// 	t.Parallel()
// 	caller := &mockWSCaller{
// 		callFn: func(_ context.Context, _ any, _ any) error {
// 			return errors.New("websocket closed")
// 		},
// 	}
// 	_, err := AccountChannels(context.Background(), caller, "rSender", "rReceiver")
// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), "account_channels")
// }

// // ─── AccountInfo ──────────────────────────────────────────────────────────────

// func TestAccountInfo_HappyPath(t *testing.T) {
// 	t.Parallel()
// 	caller := &mockWSCaller{
// 		callFn: func(_ context.Context, req, res any) error {
// 			r := req.(*AccountInfoRequest)
// 			assert.Equal(t, "account_info", r.Command)
// 			assert.Equal(t, "rAddress", r.Account)
// 			p := res.(*ResponseAccountInfo)
// 			p.Status = "success"
// 			p.Result.AccountData.Account = "rAddress"
// 			p.Result.AccountData.Balance = "1000000"
// 			return nil
// 		},
// 	}
// 	res, err := AccountInfo(context.Background(), caller, "rAddress")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// 	assert.Equal(t, "success", res.Status)
// 	assert.Equal(t, "rAddress", res.Result.AccountData.Account)
// 	assert.Equal(t, "1000000", res.Result.AccountData.Balance)
// }

// func TestAccountInfo_Error(t *testing.T) {
// 	t.Parallel()
// 	caller := &mockWSCaller{
// 		callFn: func(_ context.Context, _ any, _ any) error {
// 			return errors.New("connection reset")
// 		},
// 	}
// 	_, err := AccountInfo(context.Background(), caller, "rAddress")
// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), "account_info")
// }
