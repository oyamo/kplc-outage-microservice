package scrapper

import "github.com/oyamo/kplc-outage-microservice/pkg/model"

type MDBRepo interface {
	AddUrl(link model.Url) error
	UpdateLink(link model.Url) error
	GetUnCrawledUrl() ([]model.Url, error)
	SaveBlackoutResult(blackouts model.Blackouts) error
}
