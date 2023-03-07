package worker

import (
	"github.com/oyamo/kplc-outage-microservice/scrapping-service/internal/local/scrapper"
)

type worker struct {
	webRepo scrapper.WebRepository
}

func (w worker) DownloadUrls() error {

}

func (w worker) DownloadPdfs() error {

}

func (w worker) Schedule() chan<- error {
	errorChan := make(chan error)

	go func() {
		err := w.DownloadUrls()
		if err != nil {
			errorChan <- err
		}
	}()

	go func() {
		err := w.DownloadPdfs()
		if err != nil {
			errorChan <- err
		}
	}()

	return errorChan
}

func NewWorker() Worker {
	return &worker{}
}
