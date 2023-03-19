package scrapper

type RmqRepo interface {
	PublishId(id string) error
}
