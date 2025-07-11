package domain

import (
	"errors"
	"time"
)

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
	ErrInvalidLink       = errors.New("invalid link")
)

type Link struct {
	Id        string
	Location  string
	CreatedAt time.Time
}
