package subscription

import "github.com/oyamo/kplc-outage-microservice/pkg/model"

type MongoRepo interface {
	Subscribe(sub model.Subscription) error
	Unsubscribe(subId string) error
	GetByEmail(email string) (*model.Subscription, error)
}
