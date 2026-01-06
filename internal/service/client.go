package service

import (
	"antizapret-admin-panel/internal/entity"
	"antizapret-admin-panel/internal/repository"
	"fmt"
	"time"
)

// ClientService — это интерфейс нашего сервиса.
// Он определяет высокоуровневые бизнес-операции.
// Хендлер будет зависеть именно от этого интерфейса.
type ClientService interface {
	ListClients() ([]entity.Client, error)
	ListClientsPaginated(page, limit int) (*entity.PaginatedClients, error)
	GetClientConfigPath(name string) (string, error)
	GetClientConfigPathByType(name, configType string) (string, error)
	CreateClient(name string, expiresIn int) (*entity.Client, error)
	DeleteClient(id int) error
	GetClientByID(id int) (*entity.Client, error)
}

// clientService — конкретная реализация сервиса.
// Содержит бизнес-логику и зависит от репозитория.
// В PHP: class ClientService implements IClientService { private IClientRepository $repo; ... }
type clientService struct {
	repo repository.ClientRepository
}

// NewClientService — конструктор для нашего сервиса.
// Он принимает *интерфейс* репозитория в качестве зависимости (Dependency Injection).
func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{
		repo: repo,
	}
}

// ListClients просто вызывает метод репозитория.
func (s *clientService) ListClients() ([]entity.Client, error) {
	return s.repo.FindAll()
}

// ListClientsPaginated вызывает метод репозитория для получения пагинированных данных.
func (s *clientService) ListClientsPaginated(page, limit int) (*entity.PaginatedClients, error) {
	return s.repo.FindAllPaginated(page, limit)
}

// GetClientConfigPath также просто делегирует вызов репозиторию.
func (s *clientService) GetClientConfigPath(name string) (string, error) {
	return s.repo.FindConfigPathByName(name)
}

// GetClientConfigPathByType делегирует вызов репозиторию с указанием типа конфига.
func (s *clientService) GetClientConfigPathByType(name, configType string) (string, error) {
	return s.repo.FindConfigPathByNameAndType(name, configType)
}

// CreateClient создает нового клиента.
func (s *clientService) CreateClient(name string, expiresIn int) (*entity.Client, error) {
	err := s.repo.Create(name, expiresIn)
	if err != nil {
		return nil, err
	}
	// После успешного создания скриптом, мы можем вернуть сущность.
	// ID и другие поля здесь не так важны, фронтенд обычно просто обновляет список.
	newClient := &entity.Client{
		ID:        -1, // ID будет пересчитан при следующем вызове ListClients
		Name:      name,
		Type:      "OpenVPN", // TODO: сделать тип динамическим
		Status:    "Active",
		CreatedAt: time.Now(),
	}
	return newClient, nil
}

// GetClientByID находит клиента по его ID.
func (s *clientService) GetClientByID(id int) (*entity.Client, error) {
	clients, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, client := range clients {
		if client.ID == id {
			return &client, nil
		}
	}

	return nil, fmt.Errorf("client with ID %d not found", id)
}

// DeleteClient находит клиента по ID и удаляет его по имени.
func (s *clientService) DeleteClient(id int) error {
	client, err := s.GetClientByID(id)
	if err != nil {
		return err // Ошибка, если клиент не найден
	}

	return s.repo.DeleteByName(client.Name)
}
