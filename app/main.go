package main

import (
	"main/modules/analytics"
	"main/modules/mailout"
	"main/modules/strandsommer"
	"time"
)

func main() {
	// Call the function
	result, ctx := strandsommer.Check()
	analytics.UpdateBookings(time.Now(), ctx)
	mailout.MailOut(result)
}
