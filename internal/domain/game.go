package domain

import (
	"bytes"
	"math/rand"
)

type Player struct {
	Name string `json:"name"`
}

const IDLength int = 7 // 3 runes 1 dash 3 runes

type ID string

type Game struct {
	ID      ID       `json:"id"`
	Players []Player `json:"players"`
}

type Games map[ID]Game

func NewID() ID {
	alphaNumericRunes := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
	}

	numAlphaNumericRunes := 36

	id := make([]byte, IDLength)
	buffer := bytes.NewBuffer(id)

	for i := range id {
		if i == 3 { // midpoint, add dash
			buffer.WriteRune('-')
		} else {
			buffer.WriteRune(alphaNumericRunes[rand.Intn(numAlphaNumericRunes)])
		}
	}
	return ID(buffer.String())
}
