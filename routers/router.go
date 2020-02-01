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

	// Register controller for error handling
	beego.ErrorController(ctrl)

	// Attach controller callback object to URL paths
	beego.Router("/", ctrl, "get:GetIndex")
	beego.Router("/:entry", ctrl, "get:GetEntry")
}
