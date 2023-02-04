package pdfutil

import (
	"bytes"
	"os/exec"
)

func (*pdfReader) readPdf(path string) (bytes.Buffer, error) {
	cmd := exec.Command("pdftotext", path, "-")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return out, err
	}
	return out, nil
}
