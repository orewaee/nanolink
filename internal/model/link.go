package model

import "time"

type Link struct {
	Id        string    `yaml:"-"`
	Location  string    `yaml:"location"`
	CreatedAt time.Time `yaml:"created_at"`
}
