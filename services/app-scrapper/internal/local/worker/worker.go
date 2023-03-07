package worker

import "time"

type Config struct {
	Intervals time.Duration
}

// Worker Downloads URLS and PDFs
type Worker interface {
	DownloadUrls() error
	DownloadPdfs() error
	Schedule() error
}
