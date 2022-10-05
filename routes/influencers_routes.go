package routes

import (
	"follooow-be/handlers"

	"github.com/labstack/echo/v4"
)

func InfluencerRoute(e *echo.Echo) {
	// all routes relates to influencers comes here
	e.GET("/influencers", handlers.ListInfluencers)
	e.POST("/influencers", handlers.AddInfluencer)
	e.GET("/influencers/:influencer_id", handlers.DetailInfluencers)
	e.PUT("/influencers/:influencer_id", handlers.UpdateInfluencer)
	e.GET("/influencers/quick-find", handlers.QuickFindInfluencers)
}
