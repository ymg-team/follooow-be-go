package main

import (
	"follooow-be/configs"
	"follooow-be/routes"

	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// run database
	configs.ConnectDB()

	// routes
	routes.InfluencerRoute(e)
	routes.NewsRoute(e)
	routes.GalleriesRoute(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})
	e.Logger.Fatal(e.Start(":20223"))
}
