// Package agency implents http client and json processing to fetch a calendar and return a report.
package agency

import (
	"encoding/json"
	"github.com/codefritz/occupancy-assist/app/modules/models"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const MARKER_FREE = "Y"

func FetchReport() models.Report {
	log.Print("Start fetching occupancy data.")
	url := os.Getenv("FEWO_URL")
	strings := toOccupancyArray(fetchJson(url))
	reportDate := time.Now()
	offset := int(reportDate.Weekday()) - 1
	start := reportDate.AddDate(0, 0, -offset)
	numOccupied := 0
	result := ""

	for ix, occcupied := range strings {
		date := start.AddDate(0, 0, ix)
		if occcupied {
			numOccupied++
		}
		result += date.Format("2006-01-02") + ": " + asString(occcupied) + "\n"
	}

	log.Print("Finished fetching occupancy data.")
	log.Println("Occupied days: ", numOccupied)
	log.Println("Report date: ", reportDate)
	log.Println("Report content: ", result)

	return models.Report{Details: result, Days: numOccupied, ReportDate: reportDate}
}

func fetchJson(url string) string {

	// Fetch the json from the url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)

}

type Availability struct {
	Availabilities []string `json:"availability"`
}

type Calendar struct {
	Key Availability `json:"cal"`
}

func toOccupancyArray(jsonStr string) []bool {
	var calendar Calendar
	err := json.Unmarshal([]byte(jsonStr), &calendar)
	if err != nil {
		log.Fatal(err)
	}
	bools := make([]bool, len(calendar.Key.Availabilities))
	for i, s := range calendar.Key.Availabilities {
		bools[i] = s != MARKER_FREE
	}
	return bools
}

func asString(s bool) string {
	if s {
		return "belegt"
	}
	return "frei"
}
