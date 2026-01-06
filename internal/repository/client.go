package repository

import (
	"antizapret-admin-panel/internal/entity"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// ClientRepository — контракт
type ClientRepository interface {
	FindAll() ([]entity.Client, error)
	FindAllPaginated(page, limit int) (*entity.PaginatedClients, error)
	FindConfigPathByName(name string) (string, error)
	FindConfigPathByNameAndType(name, configType string) (string, error)
	Create(name string, expiresIn int) error
	DeleteByName(name string) error
}

// NewClientRepository — конструктор
func NewClientRepository(openvpnClientsPath string, openvpnAntizapretPath string, clientScriptPath string) ClientRepository {
	return &fileClientRepository{
		openvpnClientsPath:    openvpnClientsPath,
		openvpnAntizapretPath: openvpnAntizapretPath,
		clientScriptPath:      clientScriptPath,
	}
}

// Глобальная регулярка
var clientNameRegex = regexp.MustCompile(`^(?:vpn|antizapret)-(.+)-\(.*\)(?:-(?:udp|tcp))?\.ovpn$`)

// Хелпер для парсинга имени
func getClientName(filename string) (string, bool) {
	matches := clientNameRegex.FindStringSubmatch(filename)
	if matches == nil {
		return "", false
	}
	return matches[1], true
}

// Реализация репозитория
type fileClientRepository struct {
	openvpnClientsPath    string
	openvpnAntizapretPath string
	clientScriptPath      string
}

// FindAll читает все файлы
func (r *fileClientRepository) FindAll() ([]entity.Client, error) {
	files, err := os.ReadDir(r.openvpnClientsPath)
	if err != nil {
		return nil, err
	}

	var clients []entity.Client
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		clientName, ok := getClientName(file.Name())
		if !ok {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			log.Printf("failed to get info: %v", err)
			continue
		}

		clients = append(clients, entity.Client{
			Name:      clientName,
			Type:      "OpenVPN",
			Status:    "Active",
			CreatedAt: fileInfo.ModTime(),
		})
	}

	// Сортировка (новые сверху)
	sort.Slice(clients, func(i, j int) bool {
		return clients[i].CreatedAt.After(clients[j].CreatedAt)
	})

	// Простановка ID
	for i := range clients {
		clients[i].ID = i + 1
	}

	return clients, nil
}

func (r *fileClientRepository) FindAllPaginated(page, limit int) (*entity.PaginatedClients, error) {
	allClients, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}

	total := len(allClients)
	start := (page - 1) * limit

	// Если страница вышла за пределы, возвращаем пустой список, но с правильным Total
	if start >= total {
		return &entity.PaginatedClients{
			Total:   total,
			Clients: []entity.Client{},
		}, nil
	}

	end := start + limit
	if end > total {
		end = total
	}

	// Возвращаем структуру без полей Page/Limit, так как ты их убрал
	return &entity.PaginatedClients{
		Total:   total,
		Clients: allClients[start:end],
	}, nil
}

// FindConfigPathByName ищет конфиг в директории VPN (обратная совместимость)
func (r *fileClientRepository) FindConfigPathByName(name string) (string, error) {
	return r.FindConfigPathByNameAndType(name, "vpn")
}

// FindConfigPathByNameAndType ищет конфиг в директории соответствующего типа
func (r *fileClientRepository) FindConfigPathByNameAndType(name, configType string) (string, error) {
	var dir string
	switch configType {
	case "antizapret":
		dir = r.openvpnAntizapretPath
	default:
		dir = r.openvpnClientsPath
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		extractedName, ok := getClientName(file.Name())
		if ok && extractedName == name {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", errors.New("config file not found for client: " + name)
}

// Create
func (r *fileClientRepository) Create(name string, expiresIn int) error {
	expiresInStr := strconv.Itoa(expiresIn)
	if expiresIn <= 0 {
		expiresInStr = "3650"
	}

	cmd := exec.Command(r.clientScriptPath, "1", name, expiresInStr)
	log.Printf("Running command: %s", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create client: %w; output: %s", err, string(output))
	}
	return nil
}

// DeleteByName
func (r *fileClientRepository) DeleteByName(name string) error {
	cmd := exec.Command(r.clientScriptPath, "2", name)
	log.Printf("Running command: %s", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete client: %w; output: %s", err, string(output))
	}
	return nil
}
