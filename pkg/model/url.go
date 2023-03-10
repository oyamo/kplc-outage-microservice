package model

import "github.com/kamva/mgm/v3"

type Url struct {
	mgm.DefaultModel
	Id       int
	Link     string
	Scrapped bool
}
