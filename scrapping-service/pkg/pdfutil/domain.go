package pdfutil

import (
	"bytes"
	"time"
)

type pdfReader struct {
}

type BlackoutResult struct {
	Regions []Region
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
	Name      string
	TimeStart time.Time
	TimeStop  time.Time
	Towns     []string
}

type PdfReader interface {
	readPdf(path string) (bytes.Buffer, error)
}
