package routes

import (
	"BlogApp/controllers"
	"BlogApp/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes (require JWT token)
	userGroup := r.Group("/api/user")
	userGroup.Use(middlewares.AuthMiddleware())
	{
		userGroup.GET("/profile", controllers.GetProfile)
		userGroup.PUT("/profile", controllers.UpdateProfile)
	}
}
