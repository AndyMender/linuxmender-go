package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"

	"linuxmender/models"
	"linuxmender/utilities"
)

// RouteController is the main endpoint controller
type RouteController struct {
	beego.Controller
	EntryMgr *models.EntryManager
	CommentMgr *models.CommentManager
	EntryRecords map[int][]*models.Entry
}

// Prepare performs an initial setup before running any other method
func (ctrl *RouteController) Prepare() {
	// Attach base page layout
	ctrl.Layout = "template/layout.html"
}

// GetIndex handles the root route (the index)
// @router /
func (ctrl *RouteController) GetIndex() {
	// Load main HTML text block into LayoutContent field
	ctrl.TplName = "pages/index.html"

	// Populate remaining fields
	ctrl.Data["Title"] = "Lands of Unix"
	ctrl.Data["EntryTitle"] = "Welcome!"
	ctrl.Data["DatePosted"] = utilities.HumanizeTime(utilities.IsoToTime("2020-02-01"))
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = ""
	ctrl.Data["ValidEntry"] = false
}

// GetEntry handles routes for individual blog entry pages
// @router /posts/:entryid
func (ctrl *RouteController) GetEntry() {
	// Get entry ID and fetch matching entry details
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	entry := ctrl.EntryMgr.GetOne(entryID)

	// Abort if the blog entry doesn't exist
	if entry == nil {
		ctrl.Abort("404")
	}

	// Get comments for blog entry
	// TODO: convert `Comment.TimePosted` to a human-readable format?
	commentRecords := ctrl.CommentMgr.GetByEntry(entryID)
	if commentRecords != nil {
		ctrl.Data["Comments"] = commentRecords
	} else {
		ctrl.Data["Comments"] = make(map[int]*models.Comment)
	}

	// Load main HTML text block into LayoutContent field
	ctrl.TplName = fmt.Sprintf("pages/%v.html", entry.ID)

	// Populate remaining fields
	ctrl.Data["Title"] = entry.Title
	ctrl.Data["EntryTitle"] = entry.Title
	ctrl.Data["DatePosted"] = utilities.HumanizeTime(entry.DatePosted)
	ctrl.Data["BlogEntries"] = ctrl.EntryRecords
	ctrl.Data["EntryID"] = entry.ID
	ctrl.Data["ValidEntry"] = true
}

// GetEntryNext redirects to the subsequent blog entry page
// @router /posts/:entryid/next
func (ctrl *RouteController) GetEntryNext() {
	// Get entry ID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Repeat current entry or redirect to next if available
	nextEntryID := entryID + 1
	if entry := ctrl.EntryMgr.GetOne(nextEntryID); entry == nil {
		ctrl.Redirect(fmt.Sprintf("/posts/%v", entryID), 307)
	}

	ctrl.Redirect(fmt.Sprintf("/posts/%v", nextEntryID), 307)
}

// GetEntryPrevious redirects to the previous blog entry page
// @router /posts/:entryid/previous
func (ctrl *RouteController) GetEntryPrevious() {
	// Get entry ID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Repeat current entry or redirect to previous if available
	previousEntryID := entryID - 1
	if entry := ctrl.EntryMgr.GetOne(previousEntryID); entry == nil {
		ctrl.Redirect(fmt.Sprintf("/posts/%v", entryID), 307)
	}

	ctrl.Redirect(fmt.Sprintf("/posts/%v", previousEntryID), 307)
}
