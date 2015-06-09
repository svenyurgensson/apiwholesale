package routes

import (
	"../controllers"
	"../system"
	"../middleware"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	m "github.com/zenazn/goji/web/middleware"
)

func Include() {
	base_url := "/" + system.ApiVersion

	goji.Get( base_url + "/ping", controllers.Ping)
	goji.Post(base_url + "/session", controllers.SessionCreate)

	admin := web.New()

	goji.Handle(base_url + "/admin/*", admin)
	admin.Use(m.SubRouter)
	admin.Use(middleware.SuperSecure)

	admin.Get("/", controllers.AdminEntry)


	restricted := web.New()
	restricted.Use(middleware.TokenAuth)

	restricted.Delete(base_url + "/session", controllers.SessionDelete)

	restricted.Get( base_url + "/orders", controllers.OrdersList)
	restricted.Post(base_url + "/orders", controllers.OrderCreate)

	restricted.Get(   base_url + "/order/:order_id", controllers.OrderGet)
	restricted.Put(   base_url + "/order/:order_id", controllers.OrderUpdate)
	restricted.Delete(base_url + "/order/:order_id", controllers.OrderDelete)

	goji.Handle("/*", restricted)




}
