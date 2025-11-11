package services

import contactRepo "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"

type contact struct {
	contactRepo contactRepo.Contact
}

func NewContact(contactRepo Contact) Contact {
	return &contact{
		contactRepo: contactRepo,
	}
}
