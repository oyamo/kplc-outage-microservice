package pdfutil

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	result    BlackoutResult
	curRegion Region
	curCounty County
	curArea   BlackOutArea
)

func numberOfPages(filename string) (int, error) {
	// Run the pdfinfo command
	cmd := exec.Command("pdfinfo", filename)
	var out bytes.Buffer
	var outErr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &outErr

	err := cmd.Run()

	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Pages:") {
			countStr := strings.TrimPrefix(line, "Pages:")
			countStr = strings.TrimSpace(countStr)
			coountInt, err := strconv.Atoi(countStr)
			if err != nil {
				return 0, err
			}
			return coountInt, nil
		}
	}
	return 0, nil
}

func oneColAlign(buffer bytes.Buffer) bytes.Buffer {
	var top bytes.Buffer
	var bottom bytes.Buffer
	var scanner = bufio.NewScanner(bytes.NewReader(buffer.Bytes()))
	spaceRegex, _ := regexp.Compile(" {2,}")

	for scanner.Scan() {
		line := scanner.Text()
		segments := spaceRegex.Split(line, -1)

		// first instance
		segmentSize := len(segments)
		if segmentSize == 1 {
			top.WriteString(segments[0])
			top.WriteByte('\n')
		} else if (segmentSize == 2) && ((segments[0] == "" &&
			segments[1] != "") || (segments[1] == "" &&
			segments[0] != "")) {
			top.WriteString(segments[0])
			top.WriteByte('\n')
			bottom.WriteString(segments[1])
			bottom.WriteByte('\n')
		} else if segmentSize == 2 && strings.HasPrefix(segments[0], "DATE:") {
			top.WriteString(segments[0])
			top.WriteByte('\n')
			top.WriteString(segments[1])
			top.WriteByte('\n')
		} else if segmentSize == 3 && strings.HasPrefix(segments[2], "TIME:") {
			top.WriteString(segments[0])
			top.WriteByte('\n')
			bottom.WriteString(segments[1])
			bottom.WriteByte('\n')
			bottom.WriteString(segments[2])
			bottom.WriteByte('\n')
		} else if segmentSize == 3 && strings.HasPrefix(segments[0], "DATE:") {
			bottom.WriteString(segments[2])
			bottom.WriteByte('\n')
			top.WriteString(segments[0])
			top.WriteByte('\n')
			top.WriteString(segments[1])
			top.WriteByte('\n')
		} else if segmentSize == 4 && strings.HasPrefix(segments[0], "DATE:") &&
			segmentSize == 4 && strings.HasPrefix(segments[2], "DATE:") {
			bottom.WriteString(segments[2])
			bottom.WriteByte('\n')
			bottom.WriteString(segments[3])
			bottom.WriteByte('\n')
			top.WriteString(segments[0])
			top.WriteByte('\n')
			top.WriteString(segments[1])
			top.WriteByte('\n')
		} else if segmentSize == 2 {
			bottom.WriteString(segments[1])
			bottom.WriteByte('\n')
			top.WriteString(segments[0])
			top.WriteByte('\n')
		}
	}

	// combine the top and bottom
	top.Write(bottom.Bytes())
	return top
}

func cleanTime(t string) (res string) {

	if strings.HasSuffix(t, " A.M.") {
		res = strings.TrimSuffix(t, " A.M.")
		res += "AM"
	}
	if strings.HasSuffix(t, " P.M.") {
		res = strings.TrimSuffix(t, " P.M.")
		res += "PM"
	}

	res = strings.Replace(res, ".", ":", 1)

	builder := strings.Builder{}
	// Clear decoding issues
	for r := range res {
		if (res[r] >= '0' && res[r] <= '9') ||
			res[r] == ':' || res[r] == 'P' ||
			res[r] == 'M' || res[r] == 'A' {
			builder.WriteByte(res[r])
		}
	}

	res = builder.String()
	return
}
func formatTime(t string) (startTime string, stopTime string, err error) {
	times := bytes.Index([]byte(t), []byte(" – "))
	if times == -1 {
		times = bytes.Index([]byte(t), []byte(" - "))
		if times == -1 {
			return "", "", errors.New("invalid time - " + t)
		}
	}

	startTime = t[:times]
	stopTime = t[times+3:]

	// trim time
	startTime = startTime[len("TIME: "):]
	startTime = cleanTime(startTime)
	stopTime = cleanTime(stopTime)
	return
}

//Useful constants

func parseRegion(line string, curRegion *Region, result *BlackoutResult) bool {

	if strings.HasSuffix(line, "REGION") {
		// We have found a new region
		region := strings.TrimSuffix(line, " REGION")

		//replace hanging counties;
		lastCounty := -1
		//reduce chances of panic
		if len(curRegion.Counties) > 0 {
			lastCounty += len(curRegion.Counties)
		}

		// check if the current county is hanging
		if lastCounty > -1 {
			if curRegion.Counties[lastCounty].Name != curCounty.Name &&
				curCounty.Name != "" && curRegion.Name != "" {
				curRegion.Counties = append(curRegion.Counties, curCounty)
				curCounty = County{}
			}
		}

		// lastly check if current region is empty and there is a county hanging
		if lastCounty == -1 && curCounty.Name != "" && curRegion.Name != "" {
			curRegion.Counties = append(curRegion.Counties, curCounty)
			curCounty = County{}
		}

		// Now add the region to list of regions
		if curRegion.Name != "" || len(curRegion.Counties) != 0 {
			result.Regions = append(result.Regions, *curRegion)

		}

		*curRegion = Region{
			Name:     region,
			Counties: []County{},
		}

		return true
	}
	return false
}

func parseCounty(line string, curCounty *County, curRegion *Region) bool {
	if strings.HasPrefix(line, "PARTS OF ") {
		// We have found a new county
		if !reflect.DeepEqual(curCounty, County{}) && curCounty.Name != "" {
			curRegion.Counties = append(curRegion.Counties, *curCounty)
		}
		*curCounty = County{
			Name:  strings.TrimPrefix(line, "PARTS OF "),
			Areas: make([]BlackOutArea, 0),
		}

		return true
	}

	return false
}

func parseArea(line string, curArea *BlackOutArea, curCounty *County) bool {
	if strings.HasPrefix(line, "AREA: ") {
		area := line[6:]
		// We have found a new blackout area
		if !reflect.DeepEqual(curArea, BlackOutArea{}) && curArea.Name != "" {
			if len(curCounty.Areas) == 0 {
				curCounty.Areas = append(curCounty.Areas, *curArea)
			} else if curCounty.Areas[len(curCounty.Areas)-1].Name != curArea.Name {
				curCounty.Areas = append(curCounty.Areas, *curArea)
			}
		}

		// Reinitialise for another entry
		*curArea = BlackOutArea{
			Name:  area,
			Towns: make([]string, 0),
		}
		return true
	}
	return false
}

func parseDate(line string, curArea *BlackOutArea) (bool, error) {
	if strings.HasPrefix(line, "DATE: ") {
		// We have found the blackout date
		dateStr := line[len(line)-10:]
		t, err := time.Parse("02.01.2006", dateStr)
		if err != nil {
			if t, err = time.Parse("02.01,2006", dateStr); err != nil {
				return true, err
			}
		}
		curArea.TimeStart = t
		curArea.TimeStop = t
		curArea.TimeStartMillis = t.UnixMilli()
		curArea.TimeStopMillis = t.UnixMilli()
		return true, nil
	}

	return false, nil
}

func parseTime(line string, curArea *BlackOutArea) (bool, error) {
	if strings.HasPrefix(line, "TIME: ") {
		// We have found the blackout time

		timeStart, timeStop, err := formatTime(line)
		if err != nil {

			return true, err
		}
		t, err := time.Parse(time.Kitchen, timeStart)
		if err != nil {
			return true, err
		}

		t1, err := time.Parse(time.Kitchen, timeStop)
		if err != nil {
			return true, err
		}

		curArea.TimeStart = curArea.TimeStart.Add(t.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)))
		curArea.TimeStop = curArea.TimeStop.Add(t1.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)))
		curArea.TimeStartMillis = curArea.TimeStart.UnixMilli()
		curArea.TimeStopMillis = curArea.TimeStop.UnixMilli()
		return true, nil
	}

	return false, nil
}

func crunchSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the index of the input of a newline followed by a
	// pound sign.
	if i := strings.Index(string(data), "\n"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}
func scanTxt(buffer bytes.Buffer) (*BlackoutResult, error) {

	scanner := bufio.NewScanner(strings.NewReader(buffer.String()))
	scanner.Split(crunchSplitFunc)
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		counter++
		if parseRegion(line, &curRegion, &result) ||
			parseCounty(line, &curCounty, &curRegion) ||
			parseArea(line, &curArea, &curCounty) {
			continue
		}

		parsedDate, err := parseDate(line, &curArea)
		if err != nil {
			return nil, err
		}

		if parsedDate {
			continue
		}

		parsedTime, err := parseTime(line, &curArea)
		if err != nil {
			return nil, err
		}

		if parsedTime {
			for scanner.Scan() {
				line = scanner.Text()
				// If we reach here, we have found a list of towns
				towns := strings.Split(line, ", ")
				curArea.Towns = append(curArea.Towns, towns...)
				if strings.Contains(line, "customers") {
					break
				}
			}
			// append area
			if curCounty.Name == "" {
				curCounty.Name = curRegion.Name
			}

			curCounty.Areas = append(curCounty.Areas, curArea)
		}
	}

	if len(result.Regions) == 0 && curRegion.Name != "" {
		curRegion.Counties = append(curRegion.Counties, curCounty)
		result.Regions = append(result.Regions, curRegion)
	} else if result.Regions[len(result.Regions)-1].Name != curRegion.Name {
		curRegion.Counties = append(curRegion.Counties, curCounty)
		result.Regions = append(result.Regions, curRegion)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

func getPdfBytes(path string, page int) (bytes.Buffer, error) {
	fpage := fmt.Sprintf("-f %d", page)
	lpage := fmt.Sprintf("-l %d", page)
	args := []string{"-layout", fpage, lpage, "\"" + path + "\"", "-"}
	cmd := exec.Command("pdftotext", args...)
	cmdString := cmd.String()
	cmd = exec.Command("sh", "-c", cmdString)
	log.Printf(cmdString)
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Start()
	if err != nil {
		return out, err
	}

	if err := cmd.Wait(); err != nil {

		return out, err
	}
	return out, nil
}

func ScanPDF(path string) (*BlackoutResult, error) {
	var mainBuffer bytes.Buffer

	numPages, err := numberOfPages(path)
	if err != nil {
		return nil, err
	}

	for i := 1; i <= numPages; i++ {
		buffer, err := getPdfBytes(path, i)
		if err != nil {
			return nil, err
		}
		rdBuffer := oneColAlign(buffer)
		mainBuffer.Write(rdBuffer.Bytes())
	}

	res, err := scanTxt(mainBuffer)
	if err != nil {
		return nil, err
	}

	// Delete the PDF
	_ = os.Remove(path)

	return res, nil
}
