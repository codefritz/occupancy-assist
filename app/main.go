package main

import (
	"main/modules/agency"
	"main/modules/analytics"
	"main/modules/mailout"
	"time"
)

func main() {
	report := agency.Check()
	analytics.UpdateBookings(time.Now(), report.Days)
	mailout.MailOut(report.Content)
}
