package routers

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"

	"linuxmender/controllers"
	"linuxmender/models"
	"linuxmender/paths"
)

func init() {
	// Initialize managers
	entryManager := &models.EntryManager{
		DBName: paths.DBPath,
	}
	scoreManager := &models.ScoreManager{
		Conn: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       1,  // use default DB
		}),
	}
	commentManager := &models.CommentManager{
		DBName: paths.DBPath,
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
