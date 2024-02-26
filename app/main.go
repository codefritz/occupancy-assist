package main

import (
	"main/modules/mailout"
	"main/modules/strandsommer"
)

func main() {
	// Call the function
	result := strandsommer.Check()
	mailout.MailOut(result)
}
