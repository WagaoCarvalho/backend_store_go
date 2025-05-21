package models

import (
	"time"

	address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
)

type User struct {
	UID        int64                          `json:"uid"`
	Username   string                         `json:"username"`
	Email      string                         `json:"email"`
	Password   string                         `json:"-"`
	Status     bool                           `json:"status"`
	Version    int                            `json:"version"`
	CreatedAt  time.Time                      `json:"created_at"`
	UpdatedAt  time.Time                      `json:"updated_at"`
	Categories []user_categories.UserCategory `json:"categories,omitempty"`
	Address    *address.Address               `json:"address,omitempty"`
	Contact    *models.Contact                `json:"contact,omitempty"`
}
