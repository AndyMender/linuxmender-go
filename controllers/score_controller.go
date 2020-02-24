package controllers

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"fmt"
	"encoding/json"
)

// ScoreController is a controller for getting "likes" and "visitors" scores
type ScoreController struct {
	beego.Controller
	RedisClient *redis.Client
}

// LikeEntry adds +1 "like" to the score for a given blog entry
func (ctrl *ScoreController) LikeEntry() {
	var requestBody map[string]interface{}

	// Get entry ID and fetch matching entry details
	json.Unmarshal(ctrl.Ctx.Input.RequestBody, &requestBody)
	fmt.Printf("Post request body: %v\n", requestBody)
	ctrl.ServeJSON()
}
