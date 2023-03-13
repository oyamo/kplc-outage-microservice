package repositories

import (
	"context"
	"fmt"
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/local/scrapper"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	LinksCollectionName    = "links"
	BlackoutCollectionName = "blackouts"
)

type mongoDbRepo struct {
	Database *qmgo.Database
}

func (m mongoDbRepo) UpdateLink(link model.Url) error {
	urlCollection := m.Database.Collection(LinksCollectionName)

	// update
	err := urlCollection.UpdateOne(context.Background(), bson.M{"link": link.Link}, link)

	if err != nil {
		return err
	}

	return nil
}

func (m mongoDbRepo) AddUrl(link model.Url) error {
	urlCollection := m.Database.Collection(LinksCollectionName)

	// check if there exists a similar link
	count, err := urlCollection.Find(context.Background(), bson.M{"link": link.Link}).Count()
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("%s exists", link.Link)
	}

	// insert document
	_, err = urlCollection.InsertOne(context.Background(), link)
	if err != nil {
		return err
	}
	return nil
}

func (m mongoDbRepo) GetUnCrawledUrl() ([]model.Url, error) {
	urlCollection := m.Database.Collection(LinksCollectionName)

	// fetch uncrawled urls
	var urlResult []model.Url
	err := urlCollection.Find(context.Background(), bson.M{"scrapped": true}).All(&urlResult)
	if err != nil {
		return nil, err
	}

	return urlResult, nil
}

func (m mongoDbRepo) SaveBlackoutResult(blackouts model.Blackouts) error {
	blackoutCollection := m.Database.Collection(BlackoutCollectionName)

	err := blackouts.CalculateHash()
	if err != nil {
		return err
	}

	// check if the blackouts exists in the database
	count, err := blackoutCollection.Find(context.Background(), bson.M{"hash": blackouts.Hash}).Count()
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("blackout already exists")
	}

	// proceed to save the data
	_, err = blackoutCollection.InsertOne(context.Background(), blackouts)
	if err != nil {
		return err
	}

	return nil
}

func NewMongoRepo(db *qmgo.Database) scrapper.MDBRepo {
	return &mongoDbRepo{
		Database: db,
	}
}
