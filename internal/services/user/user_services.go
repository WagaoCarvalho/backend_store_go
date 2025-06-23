package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/auth"
	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	repositories_contact "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repositories_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	repositories_category_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

var (
	ErrInvalidEmail           = errors.New("email inválido")
	ErrCreateUser             = errors.New("erro ao criar usuário")
	ErrCreateAddress          = errors.New("erro ao criar endereço")
	ErrCreateContact          = errors.New("erro ao criar contato")
	ErrCreateCategoryRelation = errors.New("erro ao criar relação com categoria")
	ErrGetUser                = errors.New("erro ao buscar usuário")
	ErrGetVersion             = errors.New("erro ao buscar usuário")
	ErrUpdateUser             = errors.New("erro ao atualizar usuário")
	ErrDeleteUser             = errors.New("erro ao deletar usuário")
	ErrFetchAddress           = errors.New("erro ao buscar o endereço")
	ErrUpdateAddress          = errors.New("erro ao atualizar o endereço")
)

type UserService interface {
	GetAll(ctx context.Context) ([]*models_user.User, error)
	GetByID(ctx context.Context, uid int64) (*models_user.User, error)
	GetVersionByID(ctx context.Context, uid int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models_user.User, error)
	Delete(ctx context.Context, uid int64) error
	Update(
		ctx context.Context,
		user *models_user.User,
		address *models_address.Address,
	) (*models_user.User, error)
	Create(
		ctx context.Context,
		user *models_user.User,
		categoryIDs []int64,
		address *models_address.Address,
		contact *models_contact.Contact,
	) (*models_user.User, error)
}

type userService struct {
	repo                 repositories_user.UserRepository
	addressRepo          repositories_address.AddressRepository
	contactRepo          repositories_contact.ContactRepository
	categoryRelationRepo repositories_category_user.UserCategoryRelationRepository
}

func NewUserService(
	repo repositories_user.UserRepository,
	addressRepo repositories_address.AddressRepository,
	contactRepo repositories_contact.ContactRepository,
	categoryRelationRepo repositories_category_user.UserCategoryRelationRepository,
) *userService {
	return &userService{
		repo:                 repo,
		addressRepo:          addressRepo,
		contactRepo:          contactRepo,
		categoryRelationRepo: categoryRelationRepo,
	}
}

func (s *userService) Create(
	ctx context.Context,
	user *models_user.User,
	categoryIDs []int64,
	address *models_address.Address,
	contact *models_contact.Contact,
) (*models_user.User, error) {
	if !utils_validators.IsValidEmail(user.Email) {
		return nil, ErrInvalidEmail
	}

	if user.Password != "" {
		hashed, err := auth.HashPassword(user.Password)
		if err != nil {
			return nil, fmt.Errorf("erro ao hashear senha: %w", err)
		}
		user.Password = hashed
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	if createdUser == nil {
		return nil, fmt.Errorf("usuário criado é nulo")
	}

	id := int64(createdUser.UID)

	if address != nil {
		address.UserID = &id
		_, err = s.addressRepo.Create(ctx, address)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCreateAddress, err)
		}
	}

	if contact != nil {
		contact.UserID = &id
		createdContact, err := s.contactRepo.Create(ctx, contact)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrCreateContact, err)
		}
		if createdContact != nil {
			fmt.Println("Contato criado com ID:", createdContact.ID)
		}
	}

	for _, categoryID := range categoryIDs {
		relation := &models_user_category_relations.UserCategoryRelations{
			UserID:     createdUser.UID,
			CategoryID: categoryID,
		}
		_, err = s.categoryRelationRepo.Create(ctx, relation)
		if err != nil {
			return nil, fmt.Errorf("%w ID %d: %v", ErrCreateCategoryRelation, categoryID, err)
		}
	}

	return createdUser, nil
}

func (s *userService) GetAll(ctx context.Context) ([]*models_user.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	user, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}
	return user, nil
}

func (s *userService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	version, err := s.repo.GetVersionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repositories_user.ErrUserNotFound) {
			return 0, repositories_user.ErrUserNotFound
		}
		return 0, fmt.Errorf("user: erro ao obter versão: %w", err)
	}
	return version, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetUser, err)
	}
	return user, nil
}

func (s *userService) Update(
	ctx context.Context,
	user *models_user.User,
	address *models_address.Address,
) (*models_user.User, error) {
	if !utils_validators.IsValidEmail(user.Email) {
		return nil, ErrInvalidEmail
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, repositories_user.ErrUserNotFound) {
			return nil, repositories_user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	if address != nil {
		_, err := s.addressRepo.GetByID(ctx, address.ID)
		if err != nil {
			if errors.Is(err, repositories_address.ErrAddressNotFound) {
				return nil, repositories_address.ErrAddressNotFound
			}
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}

		err = s.addressRepo.Update(ctx, address)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrUpdateAddress, err)
		}
	}

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	err := s.repo.Delete(ctx, uid)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}
