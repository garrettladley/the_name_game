package websockets

import "time"

const (
	// time allowed to write a message to the peer.
	writeWait time.Duration = 10 * time.Second
	// time allowed to read the next pong message from the peer.
	pongWait time.Duration = 60 * time.Second
	// send pings to peer with this period. must be less than pongWait.
	pingPeriod time.Duration = (pongWait * 9) / 10
	// maximum message size allowed from peer.
	maxMessageSize int64 = 512
)
