package routers

import (
	"database/sql"
	"linuxmender/controllers"
	"linuxmender/models"
	"linuxmender/paths"
	"log"
	"strconv"

	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3" // provides the "sqlite3" driver in the background
)

func init() {
	// Open database with static data
	db, err := sql.Open("sqlite3", paths.EntriesDBPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Create central route controller object
	mainController := &controllers.RouteController{
		EntryRecords: models.GetEntriesAll(db),
	}

	// Register controller for error handling
	beego.ErrorController(mainController)

	// Attach controller callback object to URL paths
	beego.Router("/", mainController, "get:GetIndex")
	beego.Router("/:entry", mainController, "get:GetEntry")
}

// NextEntry tries to get the consecutive entry ID
func NextEntry(entryID string) string {
	entryNumber, err := strconv.Atoi(entryID)
	if err != nil {
		log.Printf("Entry ID: %v could not be converted to a number.\n", entryID)
		return entryID
	}

	return strconv.Itoa(entryNumber + 1)
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

	return strconv.Itoa(entryNumber)
}

// IsValidEntry checks whether the input entry ID is valid
func IsValidEntry(entryID string) bool {
	// Only numerical entries are truly valid
	if _, err := strconv.Atoi(entryID); err == nil {
		return true
	}

	return false
}
