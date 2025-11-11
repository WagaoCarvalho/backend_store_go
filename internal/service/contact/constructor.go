package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"

type contactService struct {
	repo repo.Contact
}

func NewContactService(repo repo.Contact) Contact {
	return &contactService{
		repo: repo,
	}
}
