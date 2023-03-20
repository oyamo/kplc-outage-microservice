package usecase

import (
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"strings"
	"time"
)

const SubscriberBufferSize = 1000

type usecase struct {
	mongoRepo subscription.MongoRepo
	amqpRepo  subscription.AmqpRepo
}

func (u *usecase) ProvisionNextJob(blackoutHash int64) error {
	subscribersCount, err := u.mongoRepo.GetSubscribersCount()
	timeNow := time.Now().UnixMilli()
	if err != nil {
		return err
	}

	futureBlackouts, err := u.mongoRepo.GetFutureBlackouts(0, SubscriberBufferSize)
	if err != nil {
		return err
	}

	for i := int64(0); i < subscribersCount; i += SubscriberBufferSize {
		subscriptions, err := u.mongoRepo.GetSubscribers(i, i+SubscriberBufferSize)
		if err != nil {
			return err
		}
		for j := 0; j < len(subscriptions); j++ {
			sub := subscriptions[j]
			areas := make([]model.BlackOutArea, 0)

			for j := 0; j < len(futureBlackouts); j++ {
				blackout := futureBlackouts[j]
				for regionsC := 0; regionsC < len(blackout.Regions); regionsC++ {
					regionNameUpper := strings.ToUpper(blackout.Regions[regionsC].Name)
					if strings.Index(regionNameUpper, sub.Region) != -1 {
						for countyC := 0; countyC < len(blackout.Regions[regionsC].Counties); countyC++ {
							countyName := strings.ToUpper(blackout.Regions[regionsC].Counties[countyC].Name)
							if strings.Contains(countyName, sub.County) {
								areas = append(areas, blackout.Regions[regionsC].Counties[countyC].Areas...)
							}
						}
					}
				}
			}

			err := u.mongoRepo.CreateJob(model.Job{
				Region:      sub.Region,
				County:      sub.County,
				Email:       sub.Email,
				TimeCreated: timeNow,
				TimeSent:    0,
				Sent:        false,
				Areas:       areas,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *usecase) ProvisionNextJobForUser(uuid string) error {
	//TODO implement me
	panic("implement me")
}

func (u *usecase) Subscribe(sub model.Subscription) error {
	err := u.mongoRepo.Subscribe(sub)
	if err != nil {
		return err
	}
	err = u.amqpRepo.PublishUserId(sub.UUID)
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

func NewUsecase(repo subscription.MongoRepo, amqpRepo subscription.AmqpRepo) subscription.UseCase {
	return &usecase{
		repo, amqpRepo,
	}
}
