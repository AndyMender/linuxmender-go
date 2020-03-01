package controllers

import (
	"fmt"
	"strconv"
	"encoding/json"

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

	// Get entry ID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Get user session
	sess := ctrl.GetSession("redis")
	if sess == nil {
		ctrl.SetSession("redis", uuidString)
	}
	fmt.Println(sess)

	fmt.Printf("Endorsed post with ID: %v\n", entryID)
	var requestBody map[string]interface{}

	// Get entry ID and fetch matching entry details
	json.Unmarshal(ctrl.Ctx.Input.RequestBody, &requestBody)
	fmt.Printf("Post request body: %v\n", requestBody)
	ctrl.ServeJSON()
}
