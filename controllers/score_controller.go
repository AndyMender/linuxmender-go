package controllers

import (
	"github.com/google/uuid"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"fmt"
	"strconv"
	"encoding/json"
)

// ScoreController is a controller for getting "likes" and "visitors" scores
type ScoreController struct {
	beego.Controller
	RedisClient *redis.Client
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
