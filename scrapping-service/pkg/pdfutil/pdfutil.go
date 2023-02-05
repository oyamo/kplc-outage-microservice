package pdfutil

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

var (
	skipLines = map[string]bool{
		"Interruption of":    true,
		"Electricity Supply": true,
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
)

//Useful constants

func (r *pdfReader) scanTxt(buffer bytes.Buffer) (*BlackoutResult, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))

	var result BlackoutResult
	var curRegion Region
	var curCounty County
	var curArea BlackOutArea

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || skipLines[line] {
			continue
		}

		if strings.HasSuffix(line, "REGION") {
			// We have found a new region
			if !reflect.DeepEqual(curRegion, Region{}) {
				result.Regions = append(result.Regions, curRegion)
			}
			curRegion = Region{
				Name:     strings.TrimSuffix(line, " REGION"),
				Counties: []County{},
			}
			continue
		}

		if strings.HasPrefix(line, "PARTS OF ") {
			// We have found a new county
			if !reflect.DeepEqual(curCounty, County{}) {
				curRegion.Counties = append(curRegion.Counties, curCounty)
			}
			curCounty = County{
				Name:  strings.TrimPrefix(line, "PARTS OF "),
				Areas: make([]BlackOutArea, 0),
			}
			continue
		}

		if strings.HasPrefix(line, "AREA: ") {
			// We have found a new blackout area
			if !reflect.DeepEqual(curArea, BlackOutArea{}) {
				curCounty.Areas = append(curCounty.Areas, curArea)
			}
			curArea = BlackOutArea{
				Name:  line[6:],
				Towns: make([]string, 0),
			}
			continue
		}

		if strings.HasPrefix(line, "DATE: ") {
			// We have found the blackout date
			dateStr := line[len(line)-10:]
			t, err := time.Parse("02.01.2006", dateStr)
			if err != nil {
				if t, err = time.Parse("02.01,2006", dateStr); err != nil {
					return nil, err
				}
			}
			curArea.TimeStart = t
			curArea.TimeStop = t
			continue
		}

		if strings.HasPrefix(line, "TIME: ") {
			log.Println(line)
			// We have found the blackout time
			timeStr := line[6:]
			timeStart := timeStr[:9]
			timeStop := timeStr[len(timeStr)-9:]
			timeStop = strings.TrimSpace(timeStop)

			// Clean Date to Kitchen time

			// Remove last Dot
			timeStart = timeStart[:len(timeStart)-1]
			timeStop = timeStop[:len(timeStop)-1]

			// Remove space
			timeStart = timeStart[:4] + timeStart[5:]
			timeStart = timeStart[:5] + timeStart[6:]
			timeStop = timeStop[:4] + timeStop[5:]
			timeStop = timeStop[:5] + timeStop[6:]
			_ = timeStart
			t, err := time.Parse("3.04PM", timeStart)
			if err != nil {
				if t, err = time.Parse("3.04PM", timeStart); err != nil {
					return nil, err
				}
			}

			t1, err := time.Parse("3.04PM", timeStop)
			if err != nil {
				if t1, err = time.Parse("3.04PM", timeStart); err != nil {
					return nil, err
				}
			}
			curArea.TimeStart = curArea.TimeStart.Add(t.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)))
			curArea.TimeStop = curArea.TimeStop.Add(t1.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)))
			continue
		}

		// If we reach here, we have found a list of towns
		curArea.Towns = append(curArea.Towns, strings.Split(line, ", ")...)

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	log.Printf("%+v", result)
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
