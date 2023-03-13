package worker

import "time"

type WorkType int

const (
	WorkTypeGetURl = WorkType(iota)
	WorkTypeDownloadBlackouts
)

type Config struct {
	Intervals time.Duration
}

// Worker Downloads URLS and PDFs
type Worker interface {
	Execute(err chan error, done chan bool, workType WorkType)
}
