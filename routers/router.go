package routers

import (
	"database/sql"
	"linuxmender/controllers"
	"linuxmender/models"
	"linuxmender/paths"
	"log"

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
	beego.ErrorController(&controllers.ErrorController{})

	// Attach controller callback object to URL paths
	beego.Router("/", mainController, "get:GetIndex")
	beego.Router("/posts/:entryid:int", mainController, "get:GetEntry")
	beego.Router("/posts/:entryid:int/next", mainController, "get:GetEntryNext")
	beego.Router("/posts/:entryid:int/previous", mainController, "get:GetEntryPrevious")
	beego.Router("/posts/:entryid:int/endorse", auxController, "post:LikeEntry")
}
