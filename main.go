package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"antizapret-admin-panel/internal/api"
	"antizapret-admin-panel/internal/middleware"
	"antizapret-admin-panel/internal/repository"
	"antizapret-admin-panel/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed frontend/dist/*
var embeddedFiles embed.FS

func main() {
	// Загружаем переменные из .env файла.
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используются переменные окружения ОС")
	} else {
		log.Println("Переменные из .env файла успешно загружены")
	}

	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// --- Dependency Injection ---
	// Здесь мы "собираем" наше приложение вручную.

	// 1. Получаем конфигурацию из переменных окружения
	vpnClientsPath := os.Getenv("OPENVPN_CLIENTS_PATH")
	if vpnClientsPath == "" {
		vpnClientsPath = "mock_fs/root/antizapret/client/openvpn/vpn-udp/" // Значение по умолчанию для разработки
	}
	antizapretPath := os.Getenv("OPENVPN_ANTIZAPRET_PATH")
	if antizapretPath == "" {
		antizapretPath = "mock_fs/root/antizapret/client/openvpn/antizapret-udp/"
	}
	clientScriptPath := os.Getenv("CLIENT_SCRIPT_PATH")
	if clientScriptPath == "" {
		clientScriptPath = "./mock_fs/root/antizapret/client.sh"
	}

	log.Printf("OPENVPN_CLIENTS_PATH = %s", vpnClientsPath)
	log.Printf("OPENVPN_ANTIZAPRET_PATH = %s", antizapretPath)
	log.Printf("CLIENT_SCRIPT_PATH = %s", clientScriptPath)

	// 2. Создаем Репозиторий
	clientRepo := repository.NewClientRepository(vpnClientsPath, antizapretPath, clientScriptPath)

	// 3. Создаем Сервис, внедряя в него репозиторий
	clientService := service.NewClientService(clientRepo)

	// 4. Создаем Хендлер, внедряя в него сервис
	clientHandler := api.NewClientHandler(clientService)

	// --- API Routes ---
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/login", api.LoginHandler)
		apiGroup.GET("/check-auth", middleware.AuthMiddleware(), api.CheckAuthHandler)
		// Публичный маршрут для скачивания файла по токену
		apiGroup.GET("/download/:token", api.DownloadFileHandler)

		protected := apiGroup.Group("/clients")
		protected.Use(middleware.AuthMiddleware())
		{
			// --- Обновленные роуты ---
			protected.GET("", clientHandler.GetClients)
			protected.POST("", clientHandler.CreateClient)
			protected.GET("/:id/config", clientHandler.DownloadConfig)

			// --- Старые роуты, которые пока не трогали ---
			protected.DELETE("/:id", clientHandler.DeleteClient)
			protected.GET("/:id/qr-token", clientHandler.GenerateQRToken)
		}
	}

	// Serve frontend static files AFTER API routes
	staticFiles, err := fs.Sub(embeddedFiles, "frontend/dist")
	if err != nil {
		log.Fatal("failed to get frontend/dist subtree: ", err)
	}

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Return a 404 for API routes
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Try to find the file in the static assets
		// Trim leading slash for fs.Open
		filePath := strings.TrimPrefix(path, "/")
		file, err := staticFiles.Open(filePath)
		if err == nil {
			// File exists, check if it is a directory
			stat, statErr := file.Stat()
			file.Close()

			if statErr == nil && !stat.IsDir() {
				// It's a file, serve it
				c.FileFromFS(filePath, http.FS(staticFiles))
				return
			}
		}

		// Otherwise, serve the index.html for frontend routes (SPA fallback)
		indexFile, err := staticFiles.Open("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error: index.html not found")
			return
		}
		defer indexFile.Close()

		stat, _ := indexFile.Stat()
		fileSize := stat.Size()
		buffer := make([]byte, fileSize)
		_, err = indexFile.Read(buffer)
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error: could not read index.html")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", buffer)
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
