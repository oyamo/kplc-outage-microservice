package scrapper

const (
	LinkTypePDF  = LinkType("pdf")
	LinkTypeLead = LinkType("lead")
)

type LinkType string

type Link struct {
	Url  string
	Type LinkType
}

type WebRepository interface {
	GetLinks(page string) ([]Link, error)
	GetLinksFromLead(lead string) ([]Link, error)
}
