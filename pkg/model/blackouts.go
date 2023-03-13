package model

import (
	"github.com/mitchellh/hashstructure/v2"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"
)

type Blackouts struct {
	pdfutil.BlackoutResult
	Hash uint64 `bson:"hash"`
}

func (b *Blackouts) CalculateHash() error {
	hash, err := hashstructure.Hash(b, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	b.Hash = hash
	return nil
}
