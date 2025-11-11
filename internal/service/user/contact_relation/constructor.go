package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/contact_relation"

type userContactRelationService struct {
	relationRepo repo.UserContactRelation
}

func NewUserContactRelationService(repo repo.UserContactRelation) UserContactRelation {
	return &userContactRelationService{
		relationRepo: repo,
	}
}
