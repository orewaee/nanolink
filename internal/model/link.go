package model

type Link struct {
	Id       string `yaml:"-"`
	Location string `yaml:"location"`
}
