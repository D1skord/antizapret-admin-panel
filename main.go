package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"antizapret-admin-panel/internal/api"
	"antizapret-admin-panel/internal/middleware"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	router := gin.Default()

	// API routes
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/login", api.LoginHandler)
		apiGroup.GET("/check-auth", middleware.AuthMiddleware(), api.CheckAuthHandler)

		protected := apiGroup.Group("/clients")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("", api.GetClientsHandler)
			protected.POST("", api.CreateClientHandler)
			protected.DELETE("/:id", api.DeleteClientHandler)
		}
	}

	// Serve frontend static files AFTER API routes
	staticFiles, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatal("failed to get frontend/dist subtree: ", err)
	}

	router.NoRoute(func(c *gin.Context) {
		// Return a 404 for API routes
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		// Otherwise, serve the index.html for frontend routes
		c.FileFromFS("/", http.FS(staticFiles))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
