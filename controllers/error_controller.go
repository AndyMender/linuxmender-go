package controllers

import (
	"github.com/astaxie/beego"
)


// ErrorController is the general error code controller
type ErrorController struct {
	beego.Controller
}

// Prepare performs an initial setup before running any other method
func (ctrl *ErrorController) Prepare() {
	// Attach base page layout
	ctrl.Layout = "template/layout.html"
}

// Error404 generates route details for the 404 response page
func (ctrl *ErrorController) Error404() {
	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/notfound.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Whoopsies!"
	ctrl.Data["DatePosted"] = ""
	ctrl.Data["EntryID"] = "notfound"
	ctrl.Data["ValidEntry"] = false
}
