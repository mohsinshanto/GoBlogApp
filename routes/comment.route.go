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
	// for a single post
	commentRoutes.GET("/count/:post_id", controllers.GetCommentCount)
	commentRoutes.Use(middlewares.AuthMiddleware())
	// only Authenticated users
	{
		commentRoutes.POST("/", controllers.CreateComment)
		commentRoutes.PUT("/:id", controllers.UpdateComment)
		commentRoutes.DELETE("/:id", controllers.DeleteComment)
	}
}
