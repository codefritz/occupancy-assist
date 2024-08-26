package models

import (
	"time"
)

type Report struct {
	ReportDate time.Time
	Details    string
	Days       int
}
