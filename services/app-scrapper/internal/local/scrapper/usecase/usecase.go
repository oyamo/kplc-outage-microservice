package usecase

import (
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/local/scrapper"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"
)

type useCase struct {
	MongoRepo scrapper.MDBRepo
	WebRepo   scrapper.WebRepository
}

func (u useCase) UpdateLink(link model.Url) error {
	err := u.MongoRepo.UpdateLink(link)
	if err != nil {
		return err
	}
	return nil
}

func (u useCase) GetLinks(page string) ([]scrapper.Link, error) {
	data, err := u.WebRepo.GetLinks(page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u useCase) GetLinksFromLead(lead string) ([]scrapper.Link, error) {
	data, err := u.WebRepo.GetLinksFromLead(lead)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u useCase) GenerateTmpPDF(url string) (tempPath string, err error) {
	data, err := u.WebRepo.GenerateTmpPDF(url)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (u useCase) GetBlackoutResultFromPdf(path string) (*pdfutil.BlackoutResult, error) {
	data, err := u.WebRepo.GetBlackoutResultFromPdf(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u useCase) AddUrl(link model.Url) error {
	err := u.MongoRepo.AddUrl(link)
	if err != nil {
		return err
	}
	return nil
}

func (u useCase) GetUnCrawledUrl() ([]model.Url, error) {
	data, err := u.MongoRepo.GetUnCrawledUrl()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (u useCase) SaveBlackoutResult(blackouts model.Blackouts) error {
	err := u.MongoRepo.SaveBlackoutResult(blackouts)
	if err != nil {
		return err
	}
	return err
}

func NewUseCase(mongoRepo scrapper.MDBRepo, webRepo scrapper.WebRepository) scrapper.Usecase {
	return &useCase{
		mongoRepo,
		webRepo,
	}
}
