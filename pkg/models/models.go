package models

import "time"

type Mousetrap struct {
	Id          int64
	Name        string
	OrgName     string
	Status      bool
	LastTrigger time.Time
}

type Organisation struct {
	Id int64
	Name string
	Password string
}
