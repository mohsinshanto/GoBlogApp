package routes

import (
	"BlogApp/controllers"
	"BlogApp/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(r *gin.Engine) {
	commentRoutes := r.Group("/comments")
	// Anyone can view comments
	commentRoutes.GET("/", controllers.GetAllComments)
	commentRoutes.GET("/:id", controllers.GetComment)
	commentRoutes.GET("/count/:post_id", controllers.GetCommentCount)
	commentRoutes.GET("/post/:post_id", controllers.GetCommentsByPost) // ✅ new
	commentRoutes.Use(middlewares.AuthMiddleware())
	{
		commentRoutes.POST("/", controllers.CreateComment)
		commentRoutes.PUT("/:id", controllers.UpdateComment)
		commentRoutes.DELETE("/:id", controllers.DeleteComment)
	}
}
