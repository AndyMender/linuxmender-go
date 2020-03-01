package controllers

import (
	"fmt"
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
// @router /posts/:entryid/endorse
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
// @router /posts/:entryid/likes
func (ctrl *ScoreController) GetLikes() {
	// Get entryID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	fmt.Println(entryID)
	_ = ctrl.Mgr.GetLikes(entryID)

	ctrl.ServeJSON()
}
