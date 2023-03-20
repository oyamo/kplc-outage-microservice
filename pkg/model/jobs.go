package model

import (
	"github.com/qiniu/qmgo/field"
)

const JobCollection = "jobs"

type Job struct {
	field.DefaultField `bson:",inline"`
	Region             string         `bson:"region"`
	County             string         `bson:"county"`
	Email              string         `bson:"email"`
	TimeCreated        int64          `bson:"timeCreated"`
	TimeSent           int64          `bson:"timeSent"`
	Sent               bool           `bson:"sent"`
	Areas              []BlackOutArea `bson:"area"`
}
