package utilities

import (
	"fmt"
	"log"
	"time"
)

const TimeFormat = "2006-01-02"

// IsoToTime converts an ISO7660 time string into a Time struct
func IsoToTime(timeString string) time.Time {
	// Load date string into Time object
	parsedTime, err := time.Parse(TimeFormat, timeString)
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
