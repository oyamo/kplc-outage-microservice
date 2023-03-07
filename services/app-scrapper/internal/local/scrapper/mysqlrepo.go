package scrapper

type MysqlRepo interface {
	AddUrl(link Link)
	GetUrlByID(id uint) *Link
	GetUnCrawledUrl() []Link

	// Local Cache
}
