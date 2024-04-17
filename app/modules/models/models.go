package models

import ("time")

type Report struct {
	ReportDate time.Time
	Content string
	Days    int
}