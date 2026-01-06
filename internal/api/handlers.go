package api

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

type Client struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "OpenVPN" or "WireGuard"
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Project   string    `json:"project"`
	Team      []string  `json:"team"`
	Status    string    `json:"status"`
	Budget    string    `json:"budget"`
	CreatedAt time.Time `json:"createdAt"`
}

// Mock data
var clients = []Client{
	{
		ID:        1,
		Name:      "user1",
		Type:      "OpenVPN",
		Role:      "Admin",
		Avatar:    "/images/user/user-17.jpg",
		Project:   "Antizapret Admin",
		Team:      []string{"/images/user/user-22.jpg", "/images/user/user-23.jpg"},
		Status:    "Active",
		Budget:    "5.0K",
		CreatedAt: time.Now(),
	},
	{
		ID:        2,
		Name:      "user2",
		Type:      "WireGuard",
		Role:      "User",
		Avatar:    "/images/user/user-18.jpg",
		Project:   "VPN Service",
		Team:      []string{"/images/user/user-24.jpg"},
		Status:    "Pending",
		Budget:    "2.5K",
		CreatedAt: time.Now(),
	},
}
var nextClientID = 3

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	if req.Password == adminPassword {
		c.JSON(http.StatusOK, gin.H{"token": adminPassword})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
	}
}

func CheckAuthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetClientsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, clients)
}

type CreateClientRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"` // 'openvpn' or 'wireguard'
}

func CreateClientHandler(c *gin.Context) {
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cmd *exec.Cmd
	var clientTypeDisplay string
	switch req.Type {
	case "openvpn":
		cmd = exec.Command("./mock_fs/root/antizapret/client.sh", "1", req.Name)
		clientTypeDisplay = "OpenVPN"
	case "wireguard":
		cmd = exec.Command("./mock_fs/root/antizapret/client.sh", "4", req.Name)
		clientTypeDisplay = "WireGuard"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client type. Use 'openvpn' or 'wireguard'."})
		return
	}

	log.Printf("Running command: %s", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error creating client: %v\nOutput: %s", err, string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client", "details": string(output)})
		return
	}
	log.Printf("Command output: %s", string(output))

	newClient := Client{
		ID:        nextClientID,
		Name:      req.Name,
		Type:      clientTypeDisplay,
		CreatedAt: time.Now(),
		Role:      "User",
		Avatar:    "/images/user/user-01.jpg",
		Project:   "Antizapret Admin",
		Status:    "Active",
		Budget:    "1.0K",
	}
	clients = append(clients, newClient)
	nextClientID++

	c.JSON(http.StatusCreated, newClient)
}

func DeleteClientHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	found := false
	var clientName string
	for i, client := range clients {
		if client.ID == id {
			clientName = client.Name
			clients = append(clients[:i], clients[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}

	// Placeholder for real command
	cmd := exec.Command("echo", "Deleting user", clientName)
	log.Printf("Running command: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Printf("Error deleting client: %v", err)
	}

	c.Status(http.StatusNoContent)
}
