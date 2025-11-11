package services

import (
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type userService struct {
	repo   repo.User
	hasher auth.PasswordHasher
}

func NewUserService(repo repo.User, hasher auth.PasswordHasher) User {
	return &userService{
		repo:   repo,
		hasher: hasher,
	}
}
