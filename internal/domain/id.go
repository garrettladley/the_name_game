package domain

import (
	"encoding/gob"
	"errors"
	"math/rand"

	"github.com/garrettladley/the_name_game/internal/utils"
)

const IDLength int = 7 // 3 runes 1 dash 3 runes

type ID string

func NewID() ID {
	alphaNumericRunes := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
	}

	numAlphaNumericRunes := 36

	id := make([]byte, IDLength)

	for index := range id {
		if index == 3 { // midpoint, add dash
			id[index] = '-'
		} else {
			id[index] = byte(alphaNumericRunes[rand.Intn(numAlphaNumericRunes)])
		}
	}
	return ID(string(id))
}

func (id *ID) String() string {
	return string(*id)
}

var ErrInvalidID = errors.New("invalid ID")

func ParseID(s string) (*ID, error) {
	if len(s) != IDLength {
		return nil, ErrInvalidID
	}

	for index, r := range s {
		if index == 3 {
			if r != '-' {
				return nil, ErrInvalidID
			}
		} else {
			if !utils.IsAlphanumeric(r) {
				return nil, ErrInvalidID
			}
		}
	}

	id := ID(s)
	return &id, nil
}

func (id ID) GobEncode() ([]byte, error) {
	return []byte(id), nil
}

func (id *ID) GobDecode(data []byte) error {
	*id = ID(data)
	return nil
}

func init() {
	gob.Register(ID(""))
}
