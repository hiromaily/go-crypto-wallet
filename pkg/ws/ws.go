package ws

import (
	"context"

	"github.com/pkg/errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// WS websocket object
type WS struct {
	conn *websocket.Conn
}

// New returns WS object
func New(ctx context.Context, url string) (*WS, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call websocket.Dial() %s", url)
	}

	return &WS{
		conn: conn,
	}, nil
}

// Call calls request
func (w *WS) Call(ctx context.Context, req, res interface{}) error {
	if err := wsjson.Write(ctx, w.conn, req); err != nil {
		return errors.Wrap(err, "fail to call wsjson.Write()")
	}

	if err := wsjson.Read(ctx, w.conn, res); err != nil {
		return errors.Wrap(err, "fail to call wsjson.Read()")
	}

	return nil
}

// Close disconnects
func (w *WS) Close() {
	w.conn.Close(websocket.StatusNormalClosure, "")
}
