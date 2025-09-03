package auth

import (
	"context"
	"time"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/login"
	pass "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	logger "github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	val_contact "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/contact"

	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type LoginService interface {
	Login(ctx context.Context, credentialsDTO dto.LoginCredentialsDTO) (*dto.AuthResponseDTO, error)
}

type TokenGenerator interface {
	Generate(uid int64, email string) (string, error)
}

type loginService struct {
	userRepo   repo.UserRepository
	logger     *logger.LogAdapter
	jwtManager TokenGenerator
	hasher     pass.PasswordHasher
}

func NewLoginService(repo repo.UserRepository, logger *logger.LogAdapter, jwt TokenGenerator, hasher pass.PasswordHasher) LoginService {
	return &loginService{
		userRepo:   repo,
		logger:     logger,
		jwtManager: jwt,
		hasher:     hasher,
	}
}

func (s *loginService) Login(ctx context.Context, credentialsDTO dto.LoginCredentialsDTO) (*dto.AuthResponseDTO, error) {
	const ref = "[loginService - Login] - "

	s.logger.Info(ctx, ref+logger.LogLoginInit, map[string]any{
		"email": credentialsDTO.Email,
	})

	if !val_contact.IsValidEmail(credentialsDTO.Email) {
		s.logger.Error(ctx, err_msg.ErrEmailFormat, ref+logger.LogEmailInvalid, map[string]any{
			"email": credentialsDTO.Email,
		})
		return nil, err_msg.ErrEmailFormat
	}

	credModel := credentialsDTO.ToModel()
	user, err := s.userRepo.GetByEmail(ctx, credModel.Email)
	if err != nil {
		time.Sleep(time.Second)
		s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"email": credModel.Email})
		return nil, err_msg.ErrCredentials
	}

	if err := s.hasher.Compare(user.Password, credModel.Password); err != nil {
		s.logger.Warn(ctx, ref+logger.LogPasswordInvalid, map[string]any{
			"user_id": user.UID,
			"email":   credModel.Email,
		})
		return nil, err_msg.ErrCredentials
	}

	if !user.Status {
		s.logger.Warn(ctx, ref+logger.LogAccountDisabled, map[string]any{
			"user_id": user.UID,
			"email":   user.Email,
		})
		return nil, err_msg.ErrAccountDisabled
	}

	token, err := s.jwtManager.Generate(user.UID, user.Email)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogTokenGenerationError, map[string]any{
			"user_id": user.UID,
			"email":   user.Email,
		})
		return nil, err_msg.ErrTokenGeneration
	}

	s.logger.Info(ctx, ref+logger.LogLoginSuccess, map[string]any{
		"user_id": user.UID,
		"email":   user.Email,
	})

	return &dto.AuthResponseDTO{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}
