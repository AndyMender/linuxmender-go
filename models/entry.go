package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Entry is a definition for Blog entry objects
type Entry struct {
	Title, DatePosted string
	Tags              []string
}

// GetEntries reads the entries JSON file and returns a slice of Entry records
func GetEntries(entriesPath string) map[string]Entry {
	entryRecords := make(map[string]Entry)

	entriesText, err := ioutil.ReadFile(entriesPath)
	if err != nil {
		log.Printf("%v is invalid. Cannot load entry definitions.\n", entriesPath)
		return nil
	}

	json.Unmarshal(entriesText, &entryRecords)

	return entryRecords
}
