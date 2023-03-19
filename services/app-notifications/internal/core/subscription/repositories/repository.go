package repositories

import (
	"context"
	"fmt"
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	SubCollection = "subscriptions"
)

type mongoRepo struct {
	db *qmgo.Database
}

func (m *mongoRepo) Subscribe(sub model.Subscription) error {
	collection := m.db.Collection(SubCollection)

	// Check if user has subscribed
	var subDoc model.Subscription
	err := collection.Find(context.Background(), bson.M{"email": sub.Email}).One(&subDoc)
	if err != nil {
		return err
	}

	// Update if user exists
	if subDoc.Email == sub.Email {
		err := collection.UpdateOne(context.Background(), bson.M{"email": subDoc.Email}, bson.M{"$set": &sub})
		if err != nil {
			return err
		}
		return nil
	}

	// Insert the subscription
	_, err = collection.InsertOne(context.Background(), &sub)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoRepo) Unsubscribe(subId string) error {
	collection := m.db.Collection(SubCollection)

	//Check if user has subscribed
	count, err := collection.Find(context.Background(), bson.M{"UUID": subId}).Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("%s does not exist", subId)
	}

	// Update now that it exists
	err = collection.UpdateOne(context.Background(), bson.M{"UUID": subId}, bson.M{"$set": bson.M{"subscribed": false}})
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoRepo) GetByEmail(email string) (*model.Subscription, error) {
	collection := m.db.Collection(SubCollection)

	// Fetch thedocument
	var sub model.Subscription
	err := collection.Find(context.Background(), bson.M{"email": email}).One(&sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func NewMongoRepo(db *qmgo.Database) subscription.MongoRepo {
	return &mongoRepo{
		db: db,
	}
}
