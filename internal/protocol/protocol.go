package protocol

import (
	"fmt"

	"github.com/garrettladley/the_name_game/internal/domain"
)

var (
	VERSION       byte = 1
	IDLEN         byte = 7
	SUBMITNAMEMSG byte = 1
)

type SubmitName struct {
	PlayerID domain.ID `json:"player_id"`
	Name     string    `json:"name"`
}

// [VERSION][MSGTYPE][PLAYERID][NAME]
func (s *SubmitName) MarshallBinary() (data []byte, err error) {
	data = append(data, VERSION)
	data = append(data, SUBMITNAMEMSG)
	data = append(data, []byte(s.PlayerID)...)
	data = append(data, []byte(s.Name)...)

	return data, nil
}

func (s *SubmitName) UnmarshallBinary(bytes []byte) error {
	if bytes[0] != VERSION {
		return fmt.Errorf("version mismatch %d != %d", bytes[0], VERSION)
	}
	if bytes[1] != SUBMITNAMEMSG {
		return fmt.Errorf("msg type mismatch %d != %d", bytes[1], SUBMITNAMEMSG)
	}

	s.PlayerID = domain.ID(string(bytes[2 : 2+IDLEN]))
	s.Name = string(bytes[2+IDLEN:])

	return nil
}
