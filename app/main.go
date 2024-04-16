package main

import (
	"main/modules/analytics"
	"main/modules/mailout"
	"main/modules/strandsommer"
	"time"
)

func main() {
	report := strandsommer.Check()
	analytics.UpdateBookings(time.Now(), report.Days)
	mailout.MailOut(report.Content)
}
