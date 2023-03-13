package scrapper

import (
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"
)

type Usecase interface {
	AddUrl(link model.Url) error
	UpdateLink(link model.Url) error
	GetUnCrawledUrl() ([]model.Url, error)
	SaveBlackoutResult(blackouts model.Blackouts) error
	GetLinks(page string) ([]Link, error)
	GetLinksFromLead(lead string) ([]Link, error)
	GenerateTmpPDF(url string) (tempPath string, err error)
	GetBlackoutResultFromPdf(path string) (*pdfutil.BlackoutResult, error)
}
