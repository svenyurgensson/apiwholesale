package routes

import (
    "apiwholesale/controllers"
    "apiwholesale/system"
    "apiwholesale/middleware"
    "net/http"
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


    goji.Get( base_url + "/ping", controllers.Ping)
    goji.Post(base_url + "/session", controllers.SessionCreate)

    goji.Get("/favicon.ico", controllers.Favicon)


    admin := web.New()

    goji.Handle(base_url + "/admin/*", admin)
    goji.Get(base_url + "/admin", http.RedirectHandler(base_url + "/admin/", 301))
    admin.Use(m.SubRouter)
    admin.Use(middleware.SuperSecure)
    admin.Use(c.Handler)

    admin.Get(   "/customers",             controllers.AdminCustomersList)
    admin.Post(  "/customers",             controllers.AdminCustomerCreate)
    admin.Get(   "/customer/:customer_id", controllers.AdminCustomerView)
    admin.Put(   "/customer/:customer_id", controllers.AdminCustomerUpdate)
    admin.Delete("/customer/:customer_id", controllers.AdminCustomerDelete)

    admin.Get(   "/orders",          controllers.AdminOrdersList)
    admin.Get(   "/order/:order_id", controllers.AdminOrderView)
    admin.Put(   "/order/:order_id", controllers.AdminOrderUpdate)
    admin.Delete("/order/:order_id", controllers.AdminOrderDelete)

    admin.Use(middleware.Static("public", middleware.StaticOptions{SkipLogging: true}))


    restricted := web.New()
    restricted.Use(c.Handler)
    restricted.Use(middleware.TokenAuth)

    restricted.Delete(base_url + "/session", controllers.SessionDelete)

    restricted.Get( base_url + "/orders", controllers.OrdersList)
    restricted.Post(base_url + "/orders", controllers.OrderCreate)

    restricted.Get(   base_url + "/order/:order_id", controllers.OrderGet)
    restricted.Put(   base_url + "/order/:order_id", controllers.OrderUpdate)
    restricted.Delete(base_url + "/order/:order_id", controllers.OrderDelete)

    goji.Handle("/*", restricted)

}

// Static
// https://github.com/theosomefactory/goji-static
// https://github.com/hypebeast/gojistatic
