package routes

import (
	"follooow-be/controllers"

	"github.com/labstack/echo/v4"
)

func InfluencerRoute(e *echo.Echo) {
	// all routes relates to influencers comes here
	e.GET("/influencers", controllers.ListInfluencers)
}
