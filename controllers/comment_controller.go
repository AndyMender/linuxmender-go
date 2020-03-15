package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/astaxie/beego"

	"linuxmender/models"
)

// CommentController is the main endpoint controller
type CommentController struct {
	beego.Controller
}

// SubmitComment handles POST requests with comment data
// @router /api/comments/:entryid POST
func (ctrl *CommentController) SubmitComment() {
	// Get entry ID for current entry
	entryID, _ := strconv.Atoi(ctrl.Ctx.Input.Param(":entryid"))

	// Build Comment struct
	comment := &models.Comment{}
	if err := ctrl.ParseForm(comment); err != nil {
		log.Println(err)
		ctrl.ServeJSON()
	}

	// Attach missing fields
	comment.EntryID = entryID
	comment.TimePosted = time.Now()

	fmt.Println(*comment)

	ctrl.Redirect(fmt.Sprintf("/posts/%v", entryID), 303)
}
