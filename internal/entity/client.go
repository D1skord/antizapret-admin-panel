package entity

import "time"

type Client struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type PaginatedClients struct {
	Total   int      `json:"total"`
	Clients []Client `json:"clients"`
}
