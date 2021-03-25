package models

import "time"

type Mousetrap struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	OrgName     string    `json:"org_name"`
	Status      bool      `json:"status"`
	LastTrigger time.Time `json:"last_trigger"`
}

type Organisation struct {
	Id       int64
	Name     string
	Password string
}

type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"pass"`
}
