package pdfutil

type pdfReader struct {
}

type PdfReader interface {
	ReadPdf(path string) error
}
