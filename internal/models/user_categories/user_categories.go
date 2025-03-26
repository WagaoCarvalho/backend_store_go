package models

import "time"

type UserCategory struct {
	ID          uint      `json:"id"`          // Identificador único
	Name        string    `json:"name"`        // Nome da categoria
	Description string    `json:"description"` // Descrição da categoria
	CreatedAt   time.Time `json:"created_at"`  // Data de criação
	UpdatedAt   time.Time `json:"updated_at"`  // Data da última atualização
}
