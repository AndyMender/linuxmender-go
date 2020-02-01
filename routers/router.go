package routers

import (
	"linuxmender/controllers"
	"linuxmender/models"

	"github.com/astaxie/beego"
)

func init() {
	// Create central route controller object
	ctrl := &controllers.RouteController{
		EntryRecords: models.GetEntries(),
	}

	// Attach controller callback objects to URL paths
	beego.Router("/", ctrl, "get:GetIndex")
	// TODO: separate "index" route from regular entry routes
	beego.Router("/index", ctrl, "get:GetIndex")
	beego.Router("/notfound", ctrl, "get:GetNotFound")
	beego.Router("/:entry", ctrl, "get:GetEntry")
}
