package model

type Url struct {
	Link     string `bson:"link"`
	Scrapped bool   `bson:"scrapped"`
}
