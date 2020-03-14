package routers

import (
	"database/sql"
	"log"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	_ "github.com/mattn/go-sqlite3" // provides the "sqlite3" driver in the background

	"linuxmender/controllers"
	"linuxmender/models"
	"linuxmender/paths"
)

func init() {
	// Open database with static data
	db, err := sql.Open("sqlite3", paths.EntriesDBPath)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()	// TODO: remove to avoid premature db release?

	// Initialize managers
	entryManager := &models.EntryManager{
		DB: db,
	}
	scoreManager := &models.ScoreManager{
		Conn: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       1,  // use default DB
		}),
	}

	// Create auxilliary score controller object
	auxController := &controllers.ScoreController{
		Mgr: scoreManager,
	}

	// Create central route controller object
	mainController := &controllers.RouteController{
		Mgr: entryManager,
		EntryRecords: models.EntriesByYear(entryManager.GetAll()),
	}

	// Register controller for error handling
	beego.ErrorController(&controllers.ErrorController{})

	// Attach controller callback objects to URL paths
	beego.Router("/", mainController, "get:GetIndex")
	beego.Router("/posts/:entryid:int", mainController, "get:GetEntry")
	beego.Router("/posts/:entryid:int/next", mainController, "get:GetEntryNext")
	beego.Router("/posts/:entryid:int/previous", mainController, "get:GetEntryPrevious")

	beego.Router("/api/endorse/:entryid:int", auxController, "post:LikeEntry")
	beego.Router("/api/likes/:entryid:int", auxController, "get:GetLikes")
}
