package model

import (
	"github.com/mitchellh/hashstructure/v2"
	"github.com/qiniu/qmgo/field"
	"time"
)

const (
	BlackoutsCollection = "blackouts"
)

type Blackouts struct {
	field.DefaultField `bson:",inline"`
	Regions            []Region `bson:"regions"`
	Hash               int64    `bson:"hash"`
}

type Region struct {
	Name     string
	Counties []County
}

type County struct {
	Name  string
	Areas []BlackOutArea
}

type BlackOutArea struct {
	Name            string
	TimeStart       time.Time `bson:"-"`
	TimeStartMillis int64     `bson:"timeStartMillis"`
	TimeStopMillis  int64     `bson:"timeStopMillis"`
	TimeStop        time.Time `bson:"-"`
	Towns           []string
}

func (b *Blackouts) CalculateHash() error {
	hash, err := hashstructure.Hash(b, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	b.Hash = int64(hash)
	return nil
}
