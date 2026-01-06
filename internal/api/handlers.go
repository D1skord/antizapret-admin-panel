package api

import (
	"antizapret-admin-panel/internal/service"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// --- Структуры данных запросов ---

// LoginRequest представляет структуру данных для запроса на вход.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateClientRequest представляет структуру данных для запроса на создание клиента.
type CreateClientRequest struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

// --- Управление временными токенами для скачивания ---

// DownloadTokenInfo хранит информацию о временном токене.
type DownloadTokenInfo struct {
	FilePath  string
	CreatedAt time.Time
}

const (
	TOKEN_LIFETIME         = 5 * time.Minute
	TOKEN_CLEANUP_INTERVAL = 1 * time.Minute
)

var (
	downloadTokens = make(map[string]DownloadTokenInfo)
	tokensMutex    = &sync.Mutex{}
)

func init() {
	go cleanupExpiredTokens()
}

func cleanupExpiredTokens() {
	ticker := time.NewTicker(TOKEN_CLEANUP_INTERVAL)
	defer ticker.Stop()

	for range ticker.C {
		tokensMutex.Lock()
		for token, info := range downloadTokens {
			if time.Since(info.CreatedAt) > TOKEN_LIFETIME {
				log.Printf("Удаление просроченного токена: %s", token)
				delete(downloadTokens, token)
			}
		}
		tokensMutex.Unlock()
	}
}

// generateSecureToken создает криптографически случайную строку для использования в качестве токена.
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// --- ClientHandler (Отрефакторенные обработчики клиентов) ---

// ClientHandler — это наш новый обработчик, работающий со слоем сервиса.
type ClientHandler struct {
	service service.ClientService
}

// NewClientHandler — конструктор для нашего обработчика.
func NewClientHandler(s service.ClientService) *ClientHandler {
	return &ClientHandler{service: s}
}

// GetClients обрабатывает запросы на получение списка всех клиентов.
func (h *ClientHandler) GetClients(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	clients, err := h.service.ListClientsPaginated(page, limit)
	if err != nil {
		log.Printf("Failed to retrieve clients: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve clients"})
		return
	}

	c.JSON(http.StatusOK, clients)
}

// CreateClient обрабатывает запросы на создание нового клиента.
func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != "openvpn" { // TODO: Добавить поддержку разных типов клиентов
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Создание клиентов кроме OpenVPN пока не реализовано."})
		return
	}

	newClient, err := h.service.CreateClient(req.Name, req.ExpiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newClient)
}

// DeleteClient обрабатывает запросы на удаление клиента по его ID.
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	err = h.service.DeleteClient(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client", "details": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DownloadConfig handles direct download of a client config file.
func (h *ClientHandler) DownloadConfig(c *gin.Context) {
	clientName := c.Param("id") // TODO: Сейчас мы ищем по имени, а не по ID
	configType := c.DefaultQuery("type", "vpn")

	if configType != "vpn" && configType != "antizapret" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config type. Must be 'vpn' or 'antizapret'."})
		return
	}

	filePath, err := h.service.GetClientConfigPathByType(clientName, configType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get client config path"})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Config file not found for client %s at path %s", clientName, filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration file not found."})
		return
	}

	log.Printf("Serving config file %s for client %s", filePath, clientName)
	c.FileAttachment(filePath, filepath.Base(filePath))
}

// GenerateQRToken создает временный токен для скачивания файла конфигурации.
func (h *ClientHandler) GenerateQRToken(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	configType := c.DefaultQuery("type", "vpn")
	if configType != "vpn" && configType != "antizapret" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid config type. Must be 'vpn' or 'antizapret'."})
		return
	}

	targetClient, err := h.service.GetClientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found", "details": err.Error()})
		return
	}

	if targetClient.Type != "OpenVPN" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code generation is only supported for OpenVPN clients."})
		return
	}

	filePath, err := h.service.GetClientConfigPathByType(targetClient.Name, configType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get client config path", "details": err.Error()})
		return
	}

	token, err := generateSecureToken(20)
	if err != nil {
		log.Printf("Failed to generate secure token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate download token."})
		return
	}

	tokensMutex.Lock()
	downloadTokens[token] = DownloadTokenInfo{
		FilePath:  filePath,
		CreatedAt: time.Now(),
	}
	tokensMutex.Unlock()

	log.Printf("Сгенерирован токен %s для файла %s", token, filePath)
	downloadURL := fmt.Sprintf("/api/download/%s", token)

	c.JSON(http.StatusOK, gin.H{"download_url": downloadURL})
}

// --- Прочие обработчики API (не зависят от ClientService) ---

// LoginHandler обрабатывает запросы на вход.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "password"
	}
	if req.Username == adminUsername && req.Password == adminPassword {
		c.JSON(http.StatusOK, gin.H{"token": adminPassword})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// CheckAuthHandler обрабатывает запросы на проверку аутентификации.
func CheckAuthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DownloadFileHandler обрабатывает запрос на скачивание файла по токену.
func DownloadFileHandler(c *gin.Context) {
	token := c.Param("token")
	tokensMutex.Lock()
	tokenInfo, ok := downloadTokens[token]
	if ok {
		delete(downloadTokens, token)
	}
	tokensMutex.Unlock()

	if !ok {
		log.Printf("Попытка скачивания по неверному токену: %s", token)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired token."})
		return
	}

	if time.Since(tokenInfo.CreatedAt) > TOKEN_LIFETIME {
		log.Printf("Попытка скачивания по просроченному токену: %s", token)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired token."})
		return
	}

	log.Printf("Отправка файла %s по токену %s", tokenInfo.FilePath, token)
	c.FileAttachment(tokenInfo.FilePath, filepath.Base(tokenInfo.FilePath))
}
