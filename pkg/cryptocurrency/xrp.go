package cryptocurrency

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

func NewWebSocketClient(targetURL string) (*websocket.WS, error) {
	client, err := websocket.New(context.Background(), targetURL)
	if err != nil {
		return nil, fmt.Errorf("fail to call websocket.New() for public API: %s: %w", targetURL, err)
	}

	return client, nil
}
