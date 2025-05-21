package services

import (
	"context"
	"errors"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	repositories_contact "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repositories_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	repositories_category_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type UserService interface {
	GetAll(ctx context.Context) ([]models_user.User, error)
	GetById(ctx context.Context, uid int64) (models_user.User, error)
	GetByEmail(ctx context.Context, email string) (models_user.User, error)
	Delete(ctx context.Context, uid int64) error
	Update(ctx context.Context, user *models_user.User) (models_user.User, error)
	Create(
		ctx context.Context,
		user *models_user.User,
		categoryID []int64,
		address *models_address.Address,
		contact *models_contact.Contact,
	) (models_user.User, error)
}

type userService struct {
	repo                 repositories_user.UserRepository
	addressRepo          repositories_address.AddressRepository
	contactRepo          repositories_contact.ContactRepository
	categoryRelationRepo repositories_category_user.UserCategoryRelationRepositories
}

func NewUserService(
	repo repositories_user.UserRepository,
	addressRepo repositories_address.AddressRepository,
	contactRepo repositories_contact.ContactRepository,
	categoryRelationRepo repositories_category_user.UserCategoryRelationRepositories,
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
) (models_user.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, fmt.Errorf("email inválido")
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	id := int64(createdUser.UID)
	address.UserID = &id
	_, err = s.addressRepo.Create(ctx, *address)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar endereço: %w", err)
	}

	contact.UserID = &id
	createdContact, err := s.contactRepo.Create(ctx, *contact)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar contato: %w", err)
	}

	// Exemplo: salvar ID no usuário ou logar algo
	fmt.Println("Contato criado com ID:", createdContact.ID)

	for _, categoryID := range categoryIDs {
		relation := &models_user_category_relations.UserCategoryRelations{
			UserID:     createdUser.UID,
			CategoryID: categoryID,
		}
		_, err = s.categoryRelationRepo.Create(ctx, relation)
		if err != nil {
			return models_user.User{}, fmt.Errorf("erro ao criar relação com categoria ID %d: %w", categoryID, err)
		}
	}

	return createdUser, nil
}

func (s *userService) GetAll(ctx context.Context) ([]models_user.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetById(ctx context.Context, uid int64) (models_user.User, error) {
	user, err := s.repo.GetById(ctx, uid)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, user *models_user.User) (models_user.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, fmt.Errorf("email inválido")
	}

	updatedUser, err := s.repo.Update(ctx, *user)
	if err != nil {
		if errors.Is(err, repositories_user.ErrVersionConflict) {
			return models_user.User{}, repositories_user.ErrVersionConflict
		}
		return models_user.User{}, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, uid int64) error {
	err := s.repo.Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, repositories_user.ErrRecordNotFound) {
			return fmt.Errorf("erro ao deletar usuário: %w", err)
		}
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}
