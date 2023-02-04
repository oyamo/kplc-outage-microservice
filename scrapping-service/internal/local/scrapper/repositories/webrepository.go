package repositories

import (
	"github.com/gocolly/colly"
	"github.com/oyamo/kplc-outage-microservice/scrapping-service/internal/local/scrapper"
	"log"
)

const (
	KplcBase = "https://kplc.co.ke/category/view/50/power-interruptions"
)

type webrepo struct {
}

func (w webrepo) GetLinks(page string) ([]scrapper.Link, error) {
	c := colly.NewCollector()
	c.OnHTML("div.blogSumary", func(el *colly.HTMLElement) {
		log.Println(el)
	})

	err := c.Visit(page)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nil, err
}

func NewWebRepository() scrapper.WebRepository {
	return &webrepo{}
}
