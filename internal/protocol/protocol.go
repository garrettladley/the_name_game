package protocol

type Message struct {
	Name    *string `json:"name"`
	EndGame bool    `json:"end_game"`
}
