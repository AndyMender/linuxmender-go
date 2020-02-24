package controllers

import (
	"fmt"
	"linuxmender/models"
	"strconv"

	"github.com/astaxie/beego"
)

// RouteController is the main endpoint controller
type RouteController struct {
	beego.Controller
	EntryRecords map[int]*models.Entry
}

// Prepare performs an initial setup before running any other method
func (ctrl RouteController) Prepare() {
	// Attach base page layout
	ctrl.Layout = "template/layout.html"
}

// GetIndex generates route details for the default index page
func (ctrl RouteController) GetIndex() {
	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/index.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Welcome!"
	ctrl.Data["DatePosted"] = "February 1st, 2020"
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = ""
}

// GetEntry generates route details for blog entry pages
func (ctrl RouteController) GetEntry() {
	// Additional dynamic layout sections?
	// ctrl.LayoutSections = make(map[string]string)

	// Get entry ID and fetch matching entry details
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entry"))

	entry, ok := ctrl.EntryRecords[entryID]

	// Abort if the blog entry doesn't exist
	if !ok {
		ctrl.Abort("404")
	}

	// Load main HTML text block into LayoutContent field
	ctrl.TplName = fmt.Sprintf("pages/%v.html", entryID)

	// Populate remaining fields
	ctrl.Data["Title"] = entry.Title
	ctrl.Data["EntryTitle"] = entry.Title
	ctrl.Data["DatePosted"] = entry.DatePosted
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = strconv.Itoa(entryID)
}

// Error404 generates route details for the 404 response page
func (ctrl RouteController) Error404() {
	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/notfound.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Whoopsies!"
	ctrl.Data["DatePosted"] = ""
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = "notfound"
}
