package subscription

import "github.com/oyamo/kplc-outage-microservice/pkg/model"

type MongoRepo interface {
	Subscribe(sub model.Subscription) error
	Unsubscribe(subId string) error
	GetByEmail(email string) (*model.Subscription, error)
	GetSubscribers(offset, limit int64) ([]model.Subscription, error)
	GetFutureBlackouts(offet, limit int) ([]model.Blackouts, error)
	GetBlackoutByHash(hash int64) (*model.Blackouts, error)
	CreateJob(job model.Job) error
	GetSubscribersCount() (int64, error)
}
