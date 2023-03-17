package model

import (
	"github.com/mitchellh/hashstructure/v2"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blackouts struct {
	Id primitive.ObjectID
	pdfutil.BlackoutResult
	Hash int64 `bson:"hash"`
}

func (b *Blackouts) CalculateHash() error {
	hash, err := hashstructure.Hash(b, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	b.Hash = int64(hash)
	return nil
}
