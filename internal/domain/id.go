package domain

import (
	"math/rand"
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