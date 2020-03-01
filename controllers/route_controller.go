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
func (ctrl *RouteController) Prepare() {
	// Attach base page layout
	ctrl.Layout = "template/layout.html"
}

// GetIndex generates route details for the default index page
// @router /
func (ctrl *RouteController) GetIndex() {
	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/index.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Welcome!"
	ctrl.Data["DatePosted"] = "February 1st, 2020"
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = ""
	ctrl.Data["ValidEntry"] = false
}

// GetEntry generates route details for blog entry pages
// @router /posts/:entryid
func (ctrl *RouteController) GetEntry() {
	// Get entry ID and fetch matching entry details
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	entry, ok := ctrl.EntryRecords[entryID]

	// Abort if the blog entry doesn't exist
	if !ok {
		ctrl.Abort("404")
	}

	// Load main HTML text block into LayoutContent field
	ctrl.TplName = fmt.Sprintf("pages/%v.html", entry.ID)

	// Populate remaining fields
	ctrl.Data["Title"] = entry.Title
	ctrl.Data["EntryTitle"] = entry.Title
	ctrl.Data["DatePosted"] = entry.DatePosted
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = entry.ID
	ctrl.Data["ValidEntry"] = true
}

// GetEntryNext generates entry details for the "next" entry in order
// @router /posts/:entryid/next
func (ctrl *RouteController) GetEntryNext() {
	// Get entry ID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Repeat current entry or redirect to next if available
	nextEntryID := entryID + 1
	if _, ok := ctrl.EntryRecords[nextEntryID]; !ok {
		ctrl.Redirect(fmt.Sprintf("/posts/%v", entryID), 307)
	}

	ctrl.Redirect(fmt.Sprintf("/posts/%v", nextEntryID), 307)
}

// GetEntryPrevious generates entry details for the "previous" entry in order
// @router /posts/:entryid/previous
func (ctrl *RouteController) GetEntryPrevious() {
	// Get entry ID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Repeat current entry or redirect to previous if available
	previousEntryID := entryID - 1
	if _, ok := ctrl.EntryRecords[previousEntryID]; !ok {
		ctrl.Redirect(fmt.Sprintf("/posts/%v", entryID), 307)
	}

	ctrl.Redirect(fmt.Sprintf("/posts/%v", previousEntryID), 307)
}

// Error404 generates route details for the 404 response page
func (ctrl *RouteController) Error404() {
	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/notfound.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Whoopsies!"
	ctrl.Data["DatePosted"] = ""
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = "notfound"
	ctrl.Data["ValidEntry"] = false
}
