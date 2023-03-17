package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	Id         primitive.ObjectID
	Email      string `bson:"email"`
	DeviceId   string `bson:"deviceId"`
	Region     string `bson:"region"`
	County     string `bson:"county"`
	UUID       string `bson:"UUID"`
	Subscribed bool   `bson:"subscribed"`
}
