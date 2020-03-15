package controllers

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/astaxie/beego"

	"linuxmender/models"
)

// ScoreController is a controller for getting "likes" and "visitors" scores
type ScoreController struct {
	beego.Controller
	Mgr *models.ScoreManager
}

// LikeEntry adds +1 "like" to the score for a given blog entry
// @router /api/endorse/:entryid
func (ctrl *ScoreController) LikeEntry() {
	// Generate a UUID to link session with back-end
	uuidString := uuid.New().String()

	// Get entryID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Get user session
	sessionID := ctrl.GetSession("redis")
	if sessionID == nil {
		ctrl.SetSession("redis", uuidString)
		sessionID = ctrl.GetSession("redis")
	}

	// Endorse entry
	ctrl.Mgr.LikeEntry(sessionID.(string), entryID)

	ctrl.ServeJSON()
}

// GetLikes fetches the "likes" score for a blog entry from the back-end
// @router /api/likes/:entryid
func (ctrl *ScoreController) GetLikes() {
	// Get entryID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	score := ctrl.Mgr.GetLikes(entryID)

	ctrl.Data["json"] = map[string]interface{}{
		"entryID": entryID,
		"likes": score,
	}

	ctrl.ServeJSON()
}

// AddVisit increments the visitor/hit counter for the website
// @router /api/visits POST
func (ctrl *ScoreController) AddVisit() {
	ctrl.Mgr.AddVisit()

	ctrl.ServeJSON()
}

// GetVisits fetches the "visitors" counter for the website from the back-end
// @router /api/visits GET
func (ctrl *ScoreController) GetVisits() {
	score := ctrl.Mgr.GetVisits()

	ctrl.Data["json"] = map[string]interface{}{
		"visits": score,
	}

	ctrl.ServeJSON()
}
