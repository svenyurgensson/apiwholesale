package routes

import (
	"../controllers"
	"../system"

	"github.com/zenazn/goji"
)

func Include() {
	base_url := "/" + system.ApiVersion

	goji.Get(base_url + "/ping", controllers.Ping)

	goji.Post(base_url + "/session", controllers.SessionCreate)
	goji.Delete(base_url + "/session", controllers.SessionDelete)

	goji.Get(base_url + "/orders", controllers.OrdersList)
	goji.Post(base_url + "/orders", controllers.OrderCreate)
	goji.Put(base_url + "/orders", controllers.OrderUpdate)
	goji.Delete(base_url + "/orders", controllers.OrderDelete)

}
