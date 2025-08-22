package main

import (
	"BlogApp/config"
	"BlogApp/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors" // 👈 import cors
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	proxies := os.Getenv("TRUSTED_PROXIES")
	var trustedList []string
	if proxies != "" {
		trustedList = strings.Split(proxies, ",")
	} else {
		trustedList = []string{"127.0.0.1"}
	}

	r := gin.Default()

	if err := r.SetTrustedProxies(trustedList); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// 👇 Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.RegisterBlogRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisterCommentRoutes(r)

	r.Run(":" + port)
}
