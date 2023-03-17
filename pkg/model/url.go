package model

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Url struct {
	field.DefaultField `bson:",inline"`
	Link               string `bson:"link"`
	Scrapped           bool   `bson:"scrapped"`
}

func (u *Url) DefaultId() {
	if u.Id.IsZero() {
		u.Id = primitive.NewObjectID()
	}
}
