package worker

import "github.com/oyamo/kplc-outage-microservice/scrapping-service/internal/local/scrapper"

type worker struct {
	webRepo scrapper.WebRepository
}

func (w worker) DownloadUrls() error {

}

func (w worker) DownloadPdfs() error {
  w.webRepo.
}

func (w worker) Schedule() error {
	go func() {
		err := w.DownloadUrls()
		if err != nil {

		}
	}()

	go func() {
		err := w.DownloadPdfs()
		if err != nil {

		}
	}()
}

func NewWorker() Worker {
	return &worker{}
}
