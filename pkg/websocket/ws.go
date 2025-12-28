package websocket

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// WS websocket object
type WS struct {
	conn *websocket.Conn
}

// New returns WS object
func New(ctx context.Context, url string) (*WS, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to call websocket.Dial() %s: %w", url, err)
	}

	return &WS{
		conn: conn,
	}, nil
}

// Call calls request
func (w *WS) Call(ctx context.Context, req, res any) error {
	if err := wsjson.Write(ctx, w.conn, req); err != nil {
		return fmt.Errorf("fail to call wsjson.Write(): %w", err)
	}

	if err := wsjson.Read(ctx, w.conn, res); err != nil {
		return fmt.Errorf("fail to call wsjson.Read(): %w", err)
	}

	return nil
}

// Close disconnects
func (w *WS) Close() error {
	return w.conn.Close(websocket.StatusNormalClosure, "")
}
