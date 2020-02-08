package main

import (
	"linuxmender/routers"
	_ "linuxmender/routers"

	"github.com/astaxie/beego"
)

func main() {
	// Add custom functions for in-template manipulations
	beego.AddFuncMap("nextEntry", routers.NextEntry)
	beego.AddFuncMap("previousEntry", routers.PreviousEntry)
	beego.AddFuncMap("isValidEntry", routers.IsValidEntry)

	// Run Web server
	beego.Run()
}
