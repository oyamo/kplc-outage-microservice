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
	Counties []County
}

type County struct {
	Name  string
	Areas []BlackOutArea
}

type BlackOutArea struct {
	Name string
	Time time.Time
}

type PdfReader interface {
	readPdf(path string) (bytes.Buffer, error)
}
