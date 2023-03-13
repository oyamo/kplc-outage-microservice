package model

type Url struct {
	Id       string `bson:"_id"`
	Link     string `bson:"link"`
	Scrapped bool   `bson:"scrapped"`
}
