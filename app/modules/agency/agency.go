// Package agency implents http client and json processing to featch a calendar and return a report.
package agency

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/codefritz/occupancy-assist/app/modules/models"

)

const MARKER_FREE = "Y"

func Check() models.Report {
	log.Print("Start fetching occupancy data...")
	url := os.Getenv("FEWO_URL")
	strings := toOccupancyArray(fetchJson(url))
	reportDate := time.Now()
	offset := int(reportDate.Weekday()) - 1
	start := reportDate.AddDate(0, 0, -offset)

	i := 0
	ctx := 0

	result := "*** Belegungsplan ***\n\n"

	for _, s := range strings {
		date := start.AddDate(0, 0, i)
		i++
		if s {
			ctx++
		}
		result += date.Format("2006-01-02") + ": " + occupied(s) + "\n"
	}
	log.Println(result)
	result += "Belegte Tage: " + fmt.Sprint(ctx) + "\n" + result

	return models.Report{Content: result, Days: ctx, ReportDate: reportDate}
}

func occupied(s bool) string {
	if s {
		return "belegt"
	}
	return "frei"
}

func fetchJson(url string) string {

	// Fetch the json from the url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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
