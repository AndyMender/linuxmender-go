package controllers

import (
	"fmt"
	"linuxmender/models"

	"github.com/astaxie/beego"
)

// RouteController is a composition wrapper around `beego.Controller`
type RouteController struct {
	beego.Controller
	EntryRecords map[string]models.Entry
}

// Definitions for common paths
const (
	layout = "template/layout.html"
)

// GetIndex generates route details for the default index page
// TODO: separate "index" route from regular entry routes
func (ctrl *RouteController) GetIndex() {
	// Attach base page layout
	// TODO: move call elsewhere to limit redundancy?
	ctrl.Layout = layout

	// Extract index entry
	entry, _ := ctrl.EntryRecords["index"]

	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/index.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = entry.Title
	ctrl.Data["DatePosted"] = entry.DatePosted
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
}

// GetEntry generates route details for blog entry pages
func (ctrl *RouteController) GetEntry() {
	// Attach base page layout
	// TODO: move call elsewhere to limit redundancy?
	ctrl.Layout = layout

	// Additional dynamic layout sections?
	// ctrl.LayoutSections = make(map[string]string)

	// Get entry ID and fetch matching entry details
	var entryID string = ctrl.Ctx.Input.Param(":entry")

	entry, ok := ctrl.EntryRecords[entryID]

	// Display correct entry page
	if ok {
		// Load main HTML text block into LayoutContent field
		ctrl.TplName = fmt.Sprintf("pages/%v.html", entry.Template)

		// Populate remaining fields
		ctrl.Data["Title"] = entry.Title
		ctrl.Data["EntryTitle"] = entry.Title
		ctrl.Data["DatePosted"] = entry.DatePosted
		ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	} else {
		// Default: HTTP 404 response page
		ctrl.GetNotFound()
	}
}

// GetNotFound generates route details for the 404 response page
func (ctrl *RouteController) GetNotFound() {
	// Attach base page layout
	// TODO: move call elsewhere to limit redundancy?
	ctrl.Layout = layout

	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/notfound.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Whoopsies!"
	ctrl.Data["DatePosted"] = ""
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
}
