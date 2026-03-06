package public

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // ─── ServerInfo ───────────────────────────────────────────────────────────────

// func TestServerInfo_HappyPath(t *testing.T) {
// 	t.Parallel()
// 	caller := &mockWSCaller{
// 		callFn: func(_ context.Context, req, res any) error {
// 			r := req.(*RequestCommand)
// 			assert.Equal(t, "server_info", r.Command)
// 			p := res.(*ResponseServerInfo)
// 			p.Status = "success"
// 			p.Result.Info.BuildVersion = "1.9.4"
// 			return nil
// 		},
// 	}
// 	res, err := ServerInfo(context.Background(), caller)
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// 	assert.Equal(t, "success", res.Status)
// 	assert.Equal(t, "1.9.4", res.Result.Info.BuildVersion)
// }

// func TestServerInfo_Error(t *testing.T) {
// 	t.Parallel()
// 	caller := &mockWSCaller{
// 		callFn: func(_ context.Context, _ any, _ any) error {
// 			return errors.New("node unavailable")
// 		},
// 	}
// 	_, err := ServerInfo(context.Background(), caller)
// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), "server_info")
// }
