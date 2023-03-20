package subscription

import "github.com/oyamo/kplc-outage-microservice/pkg/model"

type UseCase interface {
	Subscribe(sub model.Subscription) error
	Unsubscribe(subId string) error
	ProvisionNextJob(blackoutHash int64) error
	ProvisionNextJobForUser(uuid string) error
}
