package main

import (
	"BlogApp/config"
	"BlogApp/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Get trusted proxies from env or fallback
	proxies := os.Getenv("TRUSTED_PROXIES")
	var trustedList []string
	if proxies != "" {
		trustedList = strings.Split(proxies, ",") // support multiple IPs
	} else {
		trustedList = []string{"127.0.0.1"} // default for dev
	}
	r := gin.Default()
	// Set trusted proxies
	if err := r.SetTrustedProxies(trustedList); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}
	routes.RegisterBlogRoutes(r)
	routes.RegisterUserRoutes(r)
	routes.RegisterCommentRoutes(r)

	r.Run(":" + port)
}
