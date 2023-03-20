package repositories

import (
	"context"
	"fmt"
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type mongoRepo struct {
	db *qmgo.Database
}

func (m *mongoRepo) GetSubscribersCount() (int64, error) {
	collection := m.db.Collection(model.SubCollection)
	count, err := collection.Find(context.Background(), bson.M{}).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *mongoRepo) CreateJob(job model.Job) error {
	collection := m.db.Collection(model.JobCollection)

	_, err := collection.InsertOne(context.Background(), &job)

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoRepo) GetSubscribers(offset, limit int64) ([]model.Subscription, error) {
	collection := m.db.Collection(model.SubCollection)
	var subs []model.Subscription

	err := collection.Find(context.Background(), bson.M{"subscribed": true}).Skip(offset).Limit(limit).All(&subs)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (m *mongoRepo) GetFutureBlackouts(offet, limit int) ([]model.Blackouts, error) {
	collection := m.db.Collection(model.BlackoutsCollection)
	timeNow := time.Now().UnixMilli()

	filter := bson.M{
		"regions.counties.areas": bson.M{
			"$elemMatch": bson.M{
				"timeStartMillis": bson.M{
					"$gte": timeNow,
				},
			},
		},
	}

	var outs []model.Blackouts
	err := collection.Find(context.Background(), filter).Skip(int64(offet)).Limit(int64(limit)).All(&outs)
	if err != nil {
		return nil, err
	}

	return outs, nil
}

func (m *mongoRepo) GetBlackoutByHash(hash int64) (*model.Blackouts, error) {
	collection := m.db.Collection(model.BlackoutsCollection)
	filter := bson.M{"hash": hash}
	var out model.Blackouts

	err := collection.Find(context.Background(), filter).One(&out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (m *mongoRepo) Subscribe(sub model.Subscription) error {
	collection := m.db.Collection(model.SubCollection)

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
	collection := m.db.Collection(model.SubCollection)

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
	collection := m.db.Collection(model.SubCollection)

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
