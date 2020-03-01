package main

import (
	_ "linuxmender/routers"

	"github.com/astaxie/beego"
)

func main() {
	// Run Web server
	beego.Run()
}
