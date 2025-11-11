package services

type clientService struct {
	repo Client
}

func NewClient(repo Client) Client {
	return &clientService{
		repo: repo,
	}
}
