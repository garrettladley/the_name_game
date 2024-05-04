package protocol

import "fmt"

var (
	VERSION byte = 1
	IDLEN   byte = 7
)

type SubmitName struct {
	PlayerID string
	Name     string
}

// [VERSION][PLAYERID][NAME]
func (s *SubmitName) MarshallBinary() (data []byte, err error) {
	data = append(data, VERSION)
	data = append(data, []byte(s.PlayerID)...)
	data = append(data, []byte(s.Name)...)

	return data, nil
}

func (s *SubmitName) UnmarshallBinary(bytes []byte) error {
	if bytes[0] != VERSION {
		return fmt.Errorf("version mismatch %d != %d", bytes[0], VERSION)
	}

	s.PlayerID = string(bytes[1 : IDLEN+1])
	s.Name = string(bytes[IDLEN+1:])

	return nil
}
