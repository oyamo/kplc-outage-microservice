package scrapper

type MysqlRepo interface {
	AddUrl()
	GetUrlByID(id uint) *Link
	GetUnCrawledUrl()

	// Local Cache
}
