package pdfutil

import (
	"time"
)

type BlackoutResult struct {
	Regions []Region
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
	TimeStartMillis int64
	TimeStopMillis  int64
	TimeStop        time.Time `bson:"-"`
	Towns           []string
}
