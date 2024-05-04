package conn

import (
	"context"
	"encoding"
	"fmt"

	"nhooyr.io/websocket"
)

func Read(ctx context.Context, conn *websocket.Conn, data encoding.BinaryUnmarshaler) error {
	msgType, dataBytes, err := conn.Read(ctx)
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}
	if msgType != websocket.MessageBinary {
		return fmt.Errorf("message type is not binary")
	}
	return data.UnmarshalBinary(dataBytes)
}
