package routes

import (
	"lunch_menu/internal/handlers"
	"lunch_menu/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Admin authentication endpoints
		api.POST("/user/register", handlers.UserRegister)
		api.POST("/user/login", handlers.UserLogin)
		api.POST("/user/logout", handlers.UserLogout)

		// Restaurant endpoints
		api.GET("/restaurants", handlers.GetRestaurants)
		api.GET("/restaurants/:id", handlers.GetRestaurant)
		api.GET("/restaurants/:id/menu", handlers.GetRestaurantMenu)
		api.POST("/restaurants", middleware.AuthMiddleware, handlers.CreateRestaurant)
		api.PUT("/restaurants/:id", middleware.AuthMiddleware, handlers.UpdateRestaurant)
		api.DELETE("/restaurants/:id", middleware.AuthMiddleware, handlers.DeleteRestaurant)

		// Menu item endpoints
		api.GET("/menu-items/:id", handlers.GetMenuItem)
		api.GET("/menu-items", handlers.GetMenuItems)
		api.POST("/menu-items", middleware.AuthMiddleware, handlers.CreateMenuItem)
		api.PUT("/menu-items/:id", middleware.AuthMiddleware, handlers.UpdateMenuItem)
		api.DELETE("/menu-items/:id", middleware.AuthMiddleware, handlers.DeleteMenuItem)

		// Stats endpoint
		api.GET("/stats", handlers.GetBusinessStatistics)

		// Token refresh endpoint, do not need this endpoint as we are renewing access token in middleware itself
		// api.POST("/token/refresh", handlers.RefreshAccessToken)
	}
}
