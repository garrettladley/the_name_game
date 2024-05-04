package conn

import (
	"context"
	"encoding"
	"fmt"

	"nhooyr.io/websocket"
)

func Write(ctx context.Context, conn *websocket.Conn, data encoding.BinaryMarshaler) error {
	encodedData, err := data.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal data to binary: %w", err)
	}
	return conn.Write(ctx, websocket.MessageBinary, encodedData)
}
