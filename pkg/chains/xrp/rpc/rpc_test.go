package rpc

import (
	"context"
)

// mockWSCaller is a simple test double for WSCaller.
type mockWSCaller struct {
	callFn func(ctx context.Context, req, res any) error
}

func (m *mockWSCaller) Call(ctx context.Context, req, res any) error {
	return m.callFn(ctx, req, res)
}
