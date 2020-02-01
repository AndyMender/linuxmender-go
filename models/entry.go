package models

import (
	"encoding/json"
	"io/ioutil"
	"linuxmender/paths"
	"log"
)

// Entry is a definition for Blog entry objects
type Entry struct {
	Title, DatePosted, Template string
	Tags                        []string
}

// GetEntries reads the entries JSON file and returns a slice of Entry records
func GetEntries() map[string]Entry {
	var entryRecords map[string]Entry

	entriesText, err := ioutil.ReadFile(paths.EntriesPath)
	if err != nil {
		log.Printf("%v is invalid. Cannot load entry definitions.\n", paths.EntriesPath)
		return nil
	}

	json.Unmarshal(entriesText, &entryRecords)

	return entryRecords
}
