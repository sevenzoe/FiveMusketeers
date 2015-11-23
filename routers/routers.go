package routers

import (
	"../controllers"
	"github.com/astaxie/beego"
)

func Init() {
	// Index
	beego.Router("/", &controllers.Default{}, "GET:Index")

	// User
	beego.AutoRouter(&controllers.User{})

	// Space
	beego.Router("/space/:username", &controllers.User{}, "GET:Space")

	// Search
	beego.Router("/search", &controllers.Default{}, "GET:Search")
}
