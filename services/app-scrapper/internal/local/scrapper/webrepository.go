package scrapper

import "github.com/oyamo/kplc-outage-microservice/services/app-scrapper/pkg/pdfutil"

const (
	LinkTypePDF  = LinkType("pdf")
	LinkTypeLead = LinkType("lead")
)

type LinkType string

// Valid a check to ensure we only have two types of links
func (t LinkType) Valid() bool { return t == "pdf" || t == "lead" }

type Link struct {
	Url  string
	Type LinkType
}

type WebRepository interface {
	GetLinks(page string) ([]Link, error)
	GetLinksFromLead(lead string) ([]Link, error)
	GenerateTmpPDF(url string) (tempPath string, err error)
	GetBlackoutResultFromPdf(path string) (*pdfutil.BlackoutResult, error)
}
