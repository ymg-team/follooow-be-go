package routes

import (
	"follooow-be/handlers"

	"github.com/labstack/echo/v4"
)

func NewsRoute(e *echo.Echo) {
	// all routes relates to influencers comes here
	e.GET("/news", handlers.ListNews)
	e.GET("/news/:news_id", handlers.DetailNews)
	e.POST("/news", handlers.CreateNews)
	e.PUT("/news/:news_id", handlers.UpdateNews)
}
