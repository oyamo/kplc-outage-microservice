package worker

import (
	"github.com/oyamo/kplc-outage-microservice/pkg/model"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/core/scrapper"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/core/scrapper/repositories"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"
	"log"
	"runtime"
	"sync"
)

type worker struct {
	usecase scrapper.Usecase
}

func (w worker) downloadUrls(err chan error, done chan bool) {
	links, e := w.usecase.GetLinks(repositories.KplcBase)
	if e != nil {
		err <- e
		done <- true
		return
	}

	parsedUrls := make([]model.Url, 0)
	leadUrls := make([]scrapper.Link, 0)
	urlChan := make(chan model.Url)
	var wg sync.WaitGroup

	followLead := func(outChan chan<- model.Url, leads ...scrapper.Link) {
		defer wg.Done()
		for _, v := range leads {
			followedLinks, e := w.usecase.GetLinksFromLead(v.Url)
			if e != nil {
				err <- e
			} else {
				for _, link := range followedLinks {
					outChan <- model.Url{Link: link.Url}
				}
			}
		}
	}

	for _, v := range links {
		if v.Type == scrapper.LinkTypePDF {
			parsedUrls = append(parsedUrls, model.Url{Link: v.Url})
		} else {
			leadUrls = append(leadUrls, v)
		}
	}

	cpuLen := runtime.NumCPU()
	// If computer is single core, just create 4 worker threads

	if cpuLen == 1 {
		cpuLen = 4
	}

	for i := 0; i < cpuLen; i++ {
		wg.Add(1)
		first := (len(leadUrls) / cpuLen) * i
		last := first + (len(leadUrls) / cpuLen)
		if i == cpuLen-1 {
			last = len(leadUrls) - 1
		}

		if last == 0 {
			go followLead(urlChan, leadUrls[first:]...)
		} else {
			go followLead(urlChan, leadUrls[first:last]...)
		}
	}

	go func() {
		wg.Wait()
		close(urlChan)
	}()

	for url := range urlChan {
		parsedUrls = append(parsedUrls, url)
		// save url
	}

	for _, parsedUrl := range parsedUrls {
		e := w.usecase.AddUrl(parsedUrl)
		if e != nil {
			err <- e
		}
	}

	done <- true
}

func (w worker) parsePdfs(errChan chan error, done chan bool) {
	// Get unparsed pdfs from the db
	unparsedPdfs, err := w.usecase.GetUnCrawledUrl()
	if err != nil {
		errChan <- err
		done <- true
		return
	}

	if len(unparsedPdfs) == 0 {
		done <- true
		return
	}

	cpuLen := runtime.NumCPU()
	// If computer is single core, just create 4 worker threads

	if cpuLen == 1 {
		cpuLen = 4
	}

	var wg sync.WaitGroup
	blackOutChan := make(chan *pdfutil.BlackoutResult)

	downloadPdf := func(pdf ...model.Url) {
		defer wg.Done()
		for _, url := range pdf {
			tmpPdf, err := w.usecase.GenerateTmpPDF(url.Link)
			if err != nil {
				errChan <- err
			}

			blackoutRes, err := w.usecase.GetBlackoutResultFromPdf(tmpPdf)
			if err != nil {
				errChan <- err
			} else {
				url.Scrapped = true
				err = w.usecase.UpdateLink(url)
				if err != nil {
					errChan <- err
				}
				blackOutChan <- blackoutRes
			}
		}
	}

	for i := 0; i < cpuLen; i++ {
		wg.Add(1)
		first := (len(unparsedPdfs) / cpuLen) * i
		last := first + (len(unparsedPdfs) / cpuLen)
		if i == cpuLen-1 {
			last = len(unparsedPdfs) - 1
		}

		if last == 0 {
			go downloadPdf(unparsedPdfs[first:]...)
		} else {
			go downloadPdf(unparsedPdfs[first:last]...)
		}
	}

	go func() {
		wg.Wait()
		close(blackOutChan)
	}()

	for result := range blackOutChan {
		if result != nil {
			modelRes := model.Blackouts{
				BlackoutResult: *result,
				Hash:           0,
			}

			err = modelRes.CalculateHash()
			if err != nil {
				errChan <- err
				continue
			}

			// Save the result
			err = w.usecase.SaveBlackoutResult(modelRes)
			if err != nil {
				errChan <- err
			}
		}
	}

	done <- true
}

func (w worker) Execute(err chan error, done chan bool, workType WorkType) {
	log.Println("Running worker for work ID:", workType)
	if workType == WorkTypeGetURl {
		w.downloadUrls(err, done)
	} else {
		w.parsePdfs(err, done)
	}
	log.Println("Finished work for work ID:", workType)
}

func NewWorker(usecase scrapper.Usecase) Worker {
	return &worker{usecase}
}
