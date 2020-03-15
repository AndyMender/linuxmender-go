package utilities

import (
	"fmt"
	"log"
	"time"
)

// DateFormat is a IS07660/RFC3339 date string template
const DateFormat = "2006-01-02"

// TimeFormat is a ISO7660/RFC3339 datetime string template
const TimeFormat = "2006-01-02 03:04:05"

// IsoToTime converts an ISO7660 time string into a Time struct
func IsoToTime(timeString string) time.Time {
	// Load date string into Time object
	parsedTime, err := time.Parse(DateFormat, timeString)
	if err != nil {
		log.Printf("Couldn't properly parse date string %v. Using Epoch time.\n", timeString)
		return time.Unix(0, 0)
	}

	return parsedTime
}

// HumanizeTime converts a time.Time struct into a human-readable string
func HumanizeTime(t time.Time) string {
	// Decide on the day of month string ending 
	dayEnding := ""
	switch (t.Day() % 10) {
	case 1:
		dayEnding = "st"
	case 2:
		dayEnding = "nd"
	case 3:
		dayEnding = "rd"
	default:
		dayEnding = "th"
	}

	return fmt.Sprintf(
		"%s %d%s, %d", 
		t.Month(), 
		t.Day(), 
		dayEnding, 
		t.Year(),
	)
}
