package main

import (
	"github.com/codefritz/occupancy-assist/app/modules/agency"
	"github.com/codefritz/occupancy-assist/app/modules/analytics"
	"github.com/codefritz/occupancy-assist/app/modules/mailout"
	"time"
)

func main() {
	report := agency.Check()
	analytics.UpdateBookings(time.Now(), report.Days)
	mailout.MailOut(report.Content)
}
