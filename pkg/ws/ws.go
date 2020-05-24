package ws

import (
	"context"

	"github.com/pkg/errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// WS websocket object
type WS struct {
	ctx context.Context
	url string
}

// New returns WS object
func New(ctx context.Context, url string) *WS {
	return &WS{
		ctx: ctx,
		url: url,
	}
}

// Call calls request
func (w *WS) Call(ctx context.Context, req, res interface{}) error {
	if ctx == nil {
		ctx = w.ctx
	}
	conn, _, err := websocket.Dial(ctx, w.url, nil)
	if err != nil {
		return errors.Wrapf(err, "fail to call websocket.Dial() %s", w.url)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	err = wsjson.Write(ctx, conn, req)
	if err != nil {
		return errors.Wrap(err, "fail to call wsjson.Write()")
	}

	//var res ResponseAccountChannels
	if err = wsjson.Read(ctx, conn, res); err != nil {
		return errors.Wrap(err, "fail to call wsjson.Read()")
	}

	return nil
}
