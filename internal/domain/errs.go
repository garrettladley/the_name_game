package domain

import "errors"

var (
	ErrGameOver             = errors.New("game already over")
	ErrUserAlreadySubmitted = errors.New("user already submitted")
)
