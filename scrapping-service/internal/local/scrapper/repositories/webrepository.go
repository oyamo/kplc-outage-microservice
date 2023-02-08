package repositories

import (
	"errors"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"github.com/oyamo/kplc-outage-microservice/scrapping-service/internal/local/scrapper"
	"github.com/oyamo/kplc-outage-microservice/scrapping-service/pkg/pdfutil"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	KplcBase = "https://kplc.co.ke/category/view/50/power-interruptions"
)

type webrepo struct {
}

func (w webrepo) GetBlackoutResultFromPdf(path string) (*pdfutil.BlackoutResult, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	result, err := pdfutil.ScanPDF(path)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (w webrepo) GenerateTmpPDF(url string) (tempPath string, err error) {
	// Get a cache dir for temp storage

	if tempPath, err = os.UserCacheDir(); err != nil {
		tempPath = os.TempDir()
		if tempPath == "" {
			tempPath = "." //use current directory
		}
	}
	// create an outstream
	fileName := tempPath + "/" + uuid.NewString() + ".pdf"
	out, err := os.Create(fileName)
	defer out.Close()
	// since no special parameters are required as at the time of writing this code
	// No need of creating a http client
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	tempPath = fileName
	return
}

func (w webrepo) GetLinksFromLead(lead string) ([]scrapper.Link, error) {
	allLinks := make([]scrapper.Link, 0)
	c := colly.NewCollector()
	c.OnHTML(".attachments", func(element *colly.HTMLElement) {
		hyperLinks := element.ChildAttrs("a", "href")
		for _, v := range hyperLinks {
			if len(v) > 12 && v[len(v)-4:] == ".pdf" && v[:8] == "https://" {
				allLinks = append(allLinks, scrapper.Link{
					Url:  v,
					Type: scrapper.LinkTypePDF,
				})
			}
		}
	})

	err := c.Visit(lead)
	if err != nil {
		return nil, err
	}

	return allLinks, nil
}

func (w webrepo) GetLinks(page string) ([]scrapper.Link, error) {
	c := colly.NewCollector()
	var allLinks = make([]scrapper.Link, 0)
	c.OnHTML(".items div.blogSumary", func(el *colly.HTMLElement) {
		// Get pdf links
		var links = make([]scrapper.Link, 0)
		el.ForEach(".intro li a", func(i int, element *colly.HTMLElement) {
			var link = element.Attr("href")
			if len(link) > 12 && link[len(link)-4:] == ".pdf" && link[:8] == "https://" {
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
