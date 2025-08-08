package routes

import (
	"BlogApp/controllers"
	"BlogApp/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(r *gin.Engine) {
	posts := r.Group("/api")
	// public routes
	posts.GET("/getPosts", controllers.GetAllPosts)
	posts.GET("/singlePost/:id", controllers.GetPostById)
	// protected routes
	posts.Use(middlewares.AuthMiddleware())
	{
		posts.POST("/create", controllers.CreatePost)
		posts.PUT("/updatePost/:id", controllers.UpdateById)
		posts.DELETE("/deletePost/:id", controllers.DeleteById)
	}
}
