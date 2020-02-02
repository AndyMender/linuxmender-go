package tests

import (
	"linuxmender/models"
	"linuxmender/paths"
	"testing"
)

// TestGetEntries tests the `models.GetEntries` function
func TestGetEntries(t *testing.T) {
	// Get a valid entries JSON definition
	entryRecords := models.GetEntries(paths.EntriesPath)

	if entryRecords == nil {
		t.Error("GetEntries should've returned a map of Entry objects.")
	}

	// Get an invalid entries JSON definition
	entryRecords = models.GetEntries("invalid.json")

	if entryRecords != nil {
		t.Error("GetEntries should've returned nil.")
	}
}
