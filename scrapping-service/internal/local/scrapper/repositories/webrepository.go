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
	var allLinks = make([]scrapper.Link, 0)
	c.OnHTML(".items div.blogSumary", func(el *colly.HTMLElement) {
		// Get pdf links
		var links = make([]scrapper.Link, 0)
		el.ForEach(".intro li a", func(i int, element *colly.HTMLElement) {
			var link = element.Attr("href")
			if link != "" {
				links = append(links, scrapper.Link{
					Type: scrapper.LinkTypePDF,
					Url:  link,
				})
			}
		})

		if len(links) == 0 {
			var link = el.ChildAttr(".generictitle a", "href")
			if link != "" {
				links = append(links, scrapper.Link{
					Type: scrapper.LinkTypeLead,
					Url:  link,
				})
			}
		}
		allLinks = append(allLinks, links...)
	})

	err := c.Visit(page)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return allLinks, nil
}

func NewWebRepository() scrapper.WebRepository {
	return &webrepo{}
}
