package main

import (
	"github.com/codefritz/occupancy-assist/app/modules/agency"
	"github.com/codefritz/occupancy-assist/app/modules/analytics"
	"github.com/codefritz/occupancy-assist/app/modules/mailout"
)

func main() {
	report := agency.FetchReport()
	analytics.UpdateBookings(report.ReportDate, report.Days)
	// mailout.MailOut(report)
}
