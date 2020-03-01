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
	"github.com/go-redis/redis"
)

func init() {
	// Open database with static data
	db, err := sql.Open("sqlite3", paths.EntriesDBPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	entryManager := models.EntryManager{
		DB: db,
	}

	// Create auxilliary score controller object
	auxController := &controllers.ScoreController{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       1,  // use default DB
		}),
	}

	// Create central route controller object
	mainController := &controllers.RouteController{
		EntryRecords: entryManager.GetAll(),
	}

	// Register controller for error handling
	beego.ErrorController(mainController)

	// Attach controller callback object to URL paths
	beego.Router("/", mainController, "get:GetIndex")
	beego.Router("/posts/:entry", mainController, "get:GetEntry")
	beego.Router("/api/endorse", auxController, "post:LikeEntry")
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
