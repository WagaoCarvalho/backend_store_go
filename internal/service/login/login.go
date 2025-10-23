package auth

import (
	"context"
	"time"

	pass "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type LoginService interface {
	Login(ctx context.Context, email, password string) (*models.AuthResponse, error)
}

type TokenGenerator interface {
	Generate(uid int64, email string) (string, error)
}

type loginService struct {
	userRepo   repo.User
	jwtManager TokenGenerator
	hasher     pass.PasswordHasher
}

func NewLoginService(repo repo.User, jwt TokenGenerator, hasher pass.PasswordHasher) LoginService {
	return &loginService{
		userRepo:   repo,
		jwtManager: jwt,
		hasher:     hasher,
	}
}

func (s *loginService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	// Construir struct de credenciais
	creds := &models.LoginCredential{
		Email:    email,
		Password: password,
	}

	// Validar email e senha
	if err := creds.Validate(); err != nil {
		return nil, err // Retorna erro de validação (ex: 400)
	}

	// Buscar usuário pelo email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		time.Sleep(time.Second) // Mitigar timing attacks
		return nil, err_msg.ErrCredentials
	}

	// Comparar senha
	if err := s.hasher.Compare(user.Password, password); err != nil {
		return nil, err_msg.ErrCredentials
	}

	// Checar se conta está ativa
	if !user.Status {
		return nil, err_msg.ErrAccountDisabled
	}

	// Gerar token
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
