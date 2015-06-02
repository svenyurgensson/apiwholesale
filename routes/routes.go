package routes

import (
    "../controllers"
    "../system"
    "../middleware"

    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
)

func Include() {
    base_url := "/" + system.ApiVersion

    goji.Get(base_url + "/ping", controllers.Ping)
    goji.Post(base_url + "/session", controllers.SessionCreate)

    restricted := web.New()
    restricted.Use(middleware.TokenAuth)
    restricted.Delete(base_url + "/session", controllers.SessionDelete)
    restricted.Get(base_url + "/orders", controllers.OrdersList)
    restricted.Post(base_url + "/orders", controllers.OrderCreate)
    restricted.Put(base_url + "/order/:order_id", controllers.OrderUpdate)
    restricted.Delete(base_url + "/order/:order_id", controllers.OrderDelete)

    goji.Handle("/*", restricted)
}
