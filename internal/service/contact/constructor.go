package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"

type contact struct {
	repo repo.Contact
}

func NewContact(repo repo.Contact) Contact {
	return &contact{
		repo: repo,
	}
}
