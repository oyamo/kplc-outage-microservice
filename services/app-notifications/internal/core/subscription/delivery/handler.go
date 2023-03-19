package delivery

import (
	"context"
	"github.com/google/uuid"
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/proto/notifications"
	"github.com/oyamo/kplc-outage-microservice/services/app-notifications/internal/core/subscription"
	"log"
	"regexp"
	"strings"
)

type handler struct {
	Usecase subscription.UseCase
	notifications.UnimplementedNotificationsServer
}

func (h handler) AddSubscriber(_ context.Context, subscriber *notifications.Subscriber) (*notifications.Response, error) {
	sub := model.Subscription{
		Email:      subscriber.Email,
		DeviceId:   "", //TODO: update deviceID
		Region:     subscriber.Region,
		County:     subscriber.County,
		UUID:       uuid.New().String(),
		Subscribed: true,
	}

	var response notifications.Response
	var errStr strings.Builder
	// validate
	if sub.Email == "" {
		errStr.WriteString("Email cannot be null. ")
	}

	// Validate email
	emailRegex := regexp.MustCompile("^([a-zA-Z0-9_\\-.]+)@([a-zA-Z0-9_\\-]+)(\\.[a-zA-Z]{2,5}){1,2}$")
	if !emailRegex.MatchString(sub.Email) {
		errStr.WriteString("Invalid Email format")
	}

	// validate county
	if sub.County == "" {
		errStr.WriteString("County is nil. ")
	}

	// validate region
	if sub.Region == "" {
		errStr.WriteString("Region is nil")
	}

	if errStr.Len() != 0 {
		response.Code = 400
		response.Message = "Failed"
		response.Error = errStr.String()
		return &response, nil
	}

	err := h.Usecase.Subscribe(sub)
	if err != nil {
		response.Code = 500
		response.Message = "Failed"
		response.Error = err.Error()
		return &response, nil
	}

	response.Code = 200
	response.Message = "Success"
	return &response, nil
}

func (h handler) Unsubscribe(_ context.Context, id *notifications.SubscriptionID) (*notifications.Response, error) {
	response := notifications.Response{
		Message: "Success",
		Error:   "",
		Code:    200,
	}

	// Validate
	if id.UUID == "" {
		response.Message = "Failed"
		response.Code = 400
		return &response, nil
	}

	// Unsubscribe
	err := h.Usecase.Unsubscribe(id.UUID)
	if err != nil {
		response.Message = "Failed"
		response.Code = 400
		return &response, nil
	}

	return &response, nil
}

func (h handler) mustEmbedUnimplementedNotificationsServer() {
	log.Print("UNIMPLEMENTED")
}

func NewHandler(usecase subscription.UseCase) notifications.NotificationsServer {
	return &handler{
		Usecase: usecase,
	}
}
