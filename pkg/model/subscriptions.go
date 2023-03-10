package model

import "github.com/kamva/mgm/v3"

type Subscription struct {
	mgm.DefaultModel
	Email string
	DeviceId string
	Region string
	County string
}
