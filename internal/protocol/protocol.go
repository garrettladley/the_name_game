package protocol

import (
	"fmt"
)

var (
	VERSION       byte = 1
	IDLEN         byte = 7
	SUBMITNAMEMSG byte = 1
)

type SubmitName struct {
	Name string `json:"name"`
}

// [VERSION][MSGTYPE][NAME]
func (s *SubmitName) MarshallBinary() (data []byte, err error) {
	data = append(data, VERSION)
	data = append(data, SUBMITNAMEMSG)
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

	s.Name = string(bytes[2:])

	return nil
}
