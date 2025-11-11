package services

type contact struct {
	contactRepo Contact
}

func NewContact(contactRepo Contact) Contact {
	return &contact{
		contactRepo: contactRepo,
	}
}
