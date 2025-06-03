package models

import (
	"time"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
)

type User struct {
	UID        int64                                 `json:"uid"`
	Username   string                                `json:"username"`
	Email      string                                `json:"email"`
	Password   string                                `json:"password"`
	Status     bool                                  `json:"status"`
	Version    int                                   `json:"version"`
	CreatedAt  time.Time                             `json:"created_at"`
	UpdatedAt  time.Time                             `json:"updated_at"`
	Categories []models_user_categories.UserCategory `json:"categories,omitempty"`
	Address    *models_address.Address               `json:"address,omitempty"`
	Contact    *models_contact.Contact               `json:"contact,omitempty"`
}
