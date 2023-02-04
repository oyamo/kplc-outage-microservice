package pdfutil

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"time"
)

var (
	skipLines = map[string]bool{
		"Interruption of Electricity Supply":                                true,
		"Notice is hereby given under rule 27 of the Electric Power Rules":  true,
		"That the electricity supply will be interrupted as here under:":    true,
		"(It is necessary to interrupt supply periodically in order to":     true,
		"facilitate maintenance and upgrade of power lines to the network;": true,
		"to connect new customers or to replace power lines during road":    true,
		"construction, etc.)":                   true,
		"For further information, contact":      true,
		"The nearest Kenya Power office":        true,
		"Interruption Notices may be viewed at": true,
		"www.kplc.co.ke":                        true,
	}

	regionIndicator = "REGION"
	countyIndicator = "PARTS OF"
	dateIndicator   = "DATE:"
	timeIndicator   = "TIME:"
)

//Useful constants

func (r *pdfReader) scanTxt(buffer bytes.Buffer) (*BlackoutResult, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))

	var result BlackoutResult
	var curRegion *Region
	var curCounty *County
	var curArea *BlackOutArea

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if strings.HasSuffix(line, "REGION") {
			// We have found a new region
			curRegion = &Region{
				Name:     strings.TrimSuffix(line, " REGION"),
				Counties: []County{},
			}
			result.Region = append(result.Region, *curRegion)
			continue
		}

		if strings.HasPrefix(line, "PARTS OF ") {
			// We have found a new county
			curCounty = &County{
				Name: strings.TrimPrefix(line, "PARTS OF "),
			}
			curRegion.Counties = append(curRegion.Counties, *curCounty)
			continue
		}

		if strings.HasPrefix(line, "AREA: ") {
			// We have found a new blackout area
			curArea = &BlackOutArea{
				Name: line[6:],
			}
			curCounty.Areas = append(curCounty.Areas, *curArea)
			continue
		}

		if strings.HasPrefix(line, "DATE: ") {
			// We have found the blackout date
			dateStr := line[6:]
			t, err := time.Parse("Monday 02.01.2006", dateStr)
			if err != nil {
				return nil, err
			}
			curArea.Time = t
			continue
		}

		if strings.HasPrefix(line, "TIME: ") {
			// We have found the blackout time
			timeStr := line[6:]
			t, err := time.Parse("3:04 P.M.", timeStr)
			if err != nil {
				return nil, err
			}
			curArea.Time = curArea.Time.Add(t.Sub(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)))
			continue
		}

		// If we reach here, we have found a list of towns
		curArea.Towns = append(curArea.Towns, strings.Split(line, ", ")...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

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
