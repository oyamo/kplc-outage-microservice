package pdfutil

import (
	"bytes"
	"time"
)

type pdfReader struct {
}

type BlackoutResult struct {
	Region []Region
}

type Region struct {
	Name     string
	Counties []County
}

type County struct {
	Name  string
	Areas []BlackOutArea
}

type BlackOutArea struct {
	Name  string
	Time  time.Time
	Towns []string
}

type PdfReader interface {
	readPdf(path string) (bytes.Buffer, error)
}
