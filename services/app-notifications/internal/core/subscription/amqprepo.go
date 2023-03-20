package subscription

type AmqpRepo interface {
	PublishUserId(userId string) error
}
