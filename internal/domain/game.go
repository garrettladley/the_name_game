package domain

import (
	"math/rand"
)

type Player struct {
	Name string `json:"name"`
}

type ID string

type Game struct {
	ID      ID       `json:"id"`
	Players []Player `json:"players"`
}

type Games map[ID]Game

func NewID(length int) ID {
	alphaNumericChars := []rune{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
	}

	var result string

	for range length {
		result += string(alphaNumericChars[rand.Intn(len(alphaNumericChars))])
	}
	return ID(result)
}
