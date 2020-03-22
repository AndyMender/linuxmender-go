package routers

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"

	"linuxmender/controllers"
	"linuxmender/models"
	"linuxmender/paths"
)

func init() {
	// Initialize connection string
	connStr := fmt.Sprintf(
		"user=%v password='%v' port=%v dbname=%v sslmode=prefer", 
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	// Initialize managers
	"user=linuxmender dbname=linuxmender sslmode=prefer"
	entryManager := &models.EntryManager{
		ConnStr: connStr,
	}
	scoreManager := &models.ScoreManager{
		Conn: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       1,  // use default DB
		}),
	}
	commentManager := &models.CommentManager{
		ConnStr: connStr,
	}

	// Create auxilliary score controller object
	auxController := &controllers.ScoreController{
		Mgr: scoreManager,
	}

	commentController := &controllers.CommentController{
		Mgr: commentManager,
	}

	// Create central route controller object
	mainController := &controllers.RouteController{
		EntryMgr: entryManager,
		CommentMgr: commentManager,
	}

	// Register controller for error handling
	beego.ErrorController(&controllers.ErrorController{})

	// Attach controller callback objects to URL paths
	beego.Router("/", mainController, "get:GetIndex")
	beego.Router("/posts/:entryid:int", mainController, "get:GetEntry")
	beego.Router("/posts/:entryid:int/next", mainController, "get:GetEntryNext")
	beego.Router("/posts/:entryid:int/previous", mainController, "get:GetEntryPrevious")

	beego.Router("/api/likes/:entryid:int", auxController, "get:GetLikes")
	beego.Router("/api/likes/:entryid:int", auxController, "post:LikeEntry")
	beego.Router("/api/visits", auxController, "get:GetVisits")
	beego.Router("/api/visits", auxController, "post:AddVisit")

	beego.Router("/api/comments/:entryid:int", commentController, "post:SubmitComment")
	beego.Router("/api/comments/:entryid:int", commentController, "get:GetComments")
}
