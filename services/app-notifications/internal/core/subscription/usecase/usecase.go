package usecase

import (
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
)

type usecase struct {
	mongoRepo subscription.MongoRepo
}

func (u *usecase) Subscribe(sub model.Subscription) error {
	err := u.mongoRepo.Subscribe(sub)
	if err != nil {
		return err
	}
	return nil
}

func (u *usecase) Unsubscribe(subId string) error {
	err := u.mongoRepo.Unsubscribe(subId)
	if err != nil {
		return err
	}
	return nil
}

func NewUsecase(repo subscription.MongoRepo) subscription.UseCase {
	return &usecase{
		repo,
	}
}
