package auth

import (
	"context"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	pass "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/filter"
)

type LoginService interface {
	Login(ctx context.Context, email, password string) (*models.AuthResponse, error)
}

type TokenGenerator interface {
	Generate(uid int64, email string) (string, error)
}

type loginService struct {
	userRepo   repo.UserFilter
	jwtManager TokenGenerator
	hasher     pass.PasswordHasher
}

func NewLoginService(repo repo.UserFilter, jwt TokenGenerator, hasher pass.PasswordHasher) LoginService {
	return &loginService{
		userRepo:   repo,
		jwtManager: jwt,
		hasher:     hasher,
	}
}

func (s *loginService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	creds := &models.LoginCredential{
		Email:    email,
		Password: password,
	}

	if err := creds.Validate(); err != nil {
		return nil, err
	}

	// 🔄 Substituído GetByEmail por Filter
	userFilter := &modelFilter.UserFilter{
		Email: email,
		// Opcional: já filtrar por status ativo para evitar uma consulta extra
		// Status: utils.BoolPtr(true),
	}

	users, err := s.userRepo.Filter(ctx, userFilter)
	if err != nil {
		time.Sleep(time.Second)
		return nil, err_msg.ErrCredentials
	}

	if len(users) == 0 {
		time.Sleep(time.Second)
		return nil, err_msg.ErrCredentials
	}

	user := users[0]

	// ✅ Verificação da senha
	if err := s.hasher.Compare(user.Password, password); err != nil {
		time.Sleep(time.Second)
		return nil, err_msg.ErrCredentials
	}

	// ✅ Verificação se o usuário está ativo (já existente no seu código)
	if !user.Status {
		return nil, err_msg.ErrAccountDisabled
	}

	token, err := s.jwtManager.Generate(user.UID, user.Email)
	if err != nil {
		return nil, err_msg.ErrTokenGeneration
	}

	return &models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}
