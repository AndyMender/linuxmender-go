package routers

import (
	"fmt"
	"linuxmender/controllers"
	"linuxmender/models"
	"linuxmender/paths"
	"log"
	"strconv"

	"github.com/astaxie/beego"
)

func init() {

	// Create central route controller object
	ctrl := &controllers.RouteController{
		EntryRecords: models.GetEntries(paths.EntriesPath),
	}

	// Register controller for error handling
	beego.ErrorController(ctrl)

	// Attach controller callback object to URL paths
	beego.Router("/", ctrl, "get:GetIndex")
	beego.Router("/:entry", ctrl, "get:GetEntry")
}

// NextEntry tries to get the consecutive entry ID
func NextEntry(entryID string) string {
	entryNumber, err := strconv.Atoi(entryID)
	if err != nil {
		log.Printf("Entry ID: %v could not be converted to a number.\n", entryID)
		return entryID
	}

	return fmt.Sprintf("%d", entryNumber+1)
}

// PreviousEntry tries to get the previous entry ID
func PreviousEntry(entryID string) string {
	// Convert entryID to an integer
	entryNumber, err := strconv.Atoi(entryID)
	if err != nil {
		log.Printf("Entry ID: %v could not be converted to a number.\n", entryID)
		return entryID
	}

	entryNumber--
	// Can't move beyond first entry
	if entryNumber < 0 {
		return entryID
	}

	return fmt.Sprintf("%d", entryNumber)
}

// IsValidEntry checks whether the input entry ID is valid
func IsValidEntry(entryID string) bool {
	// Only numerical entries are truly valid
	if _, err := strconv.Atoi(entryID); err == nil {
		return true
	}

	return false
}
