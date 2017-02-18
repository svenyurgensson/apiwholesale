package routes

import (
    "apiwholesale/controllers"
    "apiwholesale/system"
    "apiwholesale/middleware"
    "github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
  m "github.com/zenazn/goji/web/middleware"
    "github.com/rs/cors"
)

func Include() {
    base_url := "/" + system.ApiVersion


    c := cors.New(cors.Options{
        AllowedOrigins: []string{},
        AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
        AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "Authorization"},
        AllowCredentials: true,
        //Debug: true,
    })
    goji.Use(c.Handler)
    goji.Use(m.RequestID)
    goji.Use(middleware.Logger)


    goji.Get( base_url + "/ping", controllers.Ping)
    goji.Post(base_url + "/session", controllers.SessionCreate)

    goji.Get("/favicon.ico", controllers.Favicon)

    goji.Get( base_url + "/me", controllers.Me)

    goji.Post( base_url + "/search", controllers.Search)

    restricted := web.New()
    restricted.Use(c.Handler)
    restricted.Use(middleware.TokenAuth)

    restricted.Delete(base_url + "/session", controllers.SessionDelete)

    restricted.Get( base_url + "/orders", controllers.OrdersList)
    restricted.Post(base_url + "/orders", controllers.OrderCreate)

    //restricted.Get(   base_url + "/order/:order_id", controllers.OrderGet)
    //restricted.Put(   base_url + "/order/:order_id", controllers.OrderUpdate)
    //restricted.Delete(base_url + "/order/:order_id", controllers.OrderDelete)

    goji.Handle("/*", restricted)

}

// Static
// https://github.com/theosomefactory/goji-static
// https://github.com/hypebeast/gojistatic
