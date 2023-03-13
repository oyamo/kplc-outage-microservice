package worker

import (
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/local/scrapper"
	"github.com/oyamo/kplc-outage-microservice/services/app-scrapper/internal/local/scrapper/worker"
	"log"
	"time"
)

func ScheduleWorker(duration time.Duration, usecase scrapper.Usecase) {
	log.Printf("Worker scheduler initialised")
	currWork := worker.WorkTypeGetURl
	for {
		log.Printf("Worker scheduler running a work")

		var doneChan = make(chan bool)
		var errChan = make(chan error)

		work := worker.NewWorker(usecase)

		go work.Execute(errChan, doneChan, currWork)

	l:
		for {
			select {
			case err, ok := <-errChan:
				if !ok {
					break l
				}
				log.Printf("error: %s", err)
			case _, ok := <-doneChan:
				{
					if !ok {
						break l
					} else {
						close(doneChan)
						close(errChan)
						break
					}
				}
			}
			if errChan == nil && doneChan == nil {
				break
			}
		}

		if currWork == worker.WorkTypeGetURl {
			currWork = worker.WorkTypeDownloadBlackouts
		} else {
			currWork = worker.WorkTypeGetURl
		}

		<-time.Tick(duration)
	}
}
