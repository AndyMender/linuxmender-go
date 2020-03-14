package utilities

import (
	"fmt"
	"log"
	"time"
)

const var TimeFormat = "2006-01-02"

// IsoToTime converts an ISO7660 time string into a Time struct
func IsoToTime(timeString String) time.Time {
	// Load date string into Time object
	parsedTime, err := time.Parse(TimeFormat, timeString)
	if err != nil {
		log.Printf("Couldn't properly parse date string %v. Using Epoch time.\n", timeString)
		return time.Unix(0, 0)
	}

	return parsedTime
}

// ParseDate converts an ISO8660 time string into a human-readable format
func ParseDate(timeString string) string {
	// Load date string into Time object
	parsedTime = IsoToTime(timeString)

	// Decide on the day of month string ending 
	dayEnding := ""
	switch (parsedTime.Day() % 10) {
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
		parsedTime.Month(), 
		parsedTime.Day(), 
		dayEnding, 
		parsedTime.Year(),
	)
}
