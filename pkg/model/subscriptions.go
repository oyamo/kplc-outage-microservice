package model

import (
	"github.com/qiniu/qmgo/field"
)

const SubCollection = "subscriptions"

type Subscription struct {
	field.DefaultField `bson:",inline"`
	Email              string `bson:"email"`
	DeviceId           string `bson:"deviceId"`
	Region             string `bson:"region"`
	County             string `bson:"county"`
	UUID               string `bson:"UUID"`
	Subscribed         bool   `bson:"subscribed"`
}
