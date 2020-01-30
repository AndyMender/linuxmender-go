package controllers

import (
	"github.com/astaxie/beego"
)

// MainController is a composition wrapper around `beego.Controller`
type MainController struct {
	beego.Controller
}

// Get method of `MainController` generates the primary HTTP response
func (c *MainController) Get() {
	// Base page layout
	c.Layout = "template/layout.html"

	// Fields loaded into the template
	c.Data["Title"] = "Lands of Unix"
	c.Data["EntryTitle"] = "Welcome!"
	c.TplName = "index.html"

	// Blog post listing
	BlogEntries := make(map[string]string)
	BlogEntries["Fishing with Ubuntu"] = "/"
	BlogEntries["I Love nVidia"] = "/"
	BlogEntries["A Very Long Blog Post Title Which I Don't Want. Does This Even Fit???"] = "/"
	c.Data["BlogEntries"] = BlogEntries

	// Additional dynamic layout sections
	c.LayoutSections = make(map[string]string)

}
