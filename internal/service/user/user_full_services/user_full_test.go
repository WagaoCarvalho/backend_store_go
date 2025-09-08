package services

import (
	"context"
	"errors"
	"testing"

	mock_tx "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo"
	mock_address "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/address"
	mock_contact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	mock_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	mock_user_full "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	model_address "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	model_user "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	model_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	model_user_full "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserFullService() (
	*mock_user_full.MockUserFullRepository,
	*mock_address.MockAddressRepository,
	*mock_contact.MockContactRepository,
	*mock_user_cat_rel.MockUserCategoryRelationRepo,
	*mock_tx.MockTx,
	*MockHasher,
	UserFullService,
) {
	mockUserRepo := new(mock_user_full.MockUserFullRepository)
	mockAddressRepo := new(mock_address.MockAddressRepository)
	mockContactRepo := new(mock_contact.MockContactRepository)
	mockRelationRepo := new(mock_user_cat_rel.MockUserCategoryRelationRepo)
	mockHasher := new(MockHasher)
	mockTx := new(mock_tx.MockTx)

	service := NewUserFullService(
		mockUserRepo,
		mockAddressRepo,
		mockContactRepo,
		mockRelationRepo,
		mockHasher,
	)

	return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, service
}

func TestUserFullService_CreateFull(t *testing.T) {
	ctx := context.Background()

	t.Run("userFull_nulo", func(t *testing.T) {
		_, _, _, _, _, _, service := setupUserFullService()

		result, err := service.CreateFull(ctx, nil)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "userFull é nulo")
	})

	t.Run("falha_quando_user_nil", func(t *testing.T) {
		_, _, _, _, _, _, service := setupUserFullService()

		invalidUser := &model_user_full.UserFull{
			Address:    &model_address.Address{Street: "Rua Teste"},
			Contact:    &model_contact.Contact{Phone: "1112345678"},
			Categories: []model_user_cat_rel.UserCategory{{ID: 1}},
		}

		result, err := service.CreateFull(ctx, invalidUser)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "userFull é nulo")
	})

	t.Run("falha_quando_address_nil", func(t *testing.T) {
		mockUserRepo, _, _, _, mockTx, mockHasher, service := setupUserFullService()

		invalidUser := &model_user_full.UserFull{
			User:       &model_user.User{Username: "user", Email: "user@test.com"},
			Contact:    &model_contact.Contact{Phone: "1112345678"},
			Categories: []model_user_cat_rel.UserCategory{{ID: 1}},
		}

		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		result, err := service.CreateFull(ctx, invalidUser)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "endereço")
		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertNotCalled(t, "Hash")
	})

	t.Run("falha_quando_contact_nil", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, _, _, mockTx, mockHasher, service := setupUserFullService()

		invalidUser := &model_user_full.UserFull{
			User:       &model_user.User{Username: "user", Email: "user@test.com"},
			Address:    &model_address.Address{Street: "Rua Teste"},
			Categories: []model_user_cat_rel.UserCategory{{ID: 1}},
		}

		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		result, err := service.CreateFull(ctx, invalidUser)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "contato")
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockHasher.AssertNotCalled(t, "Hash")
	})

	t.Run("falha_quando_sem_categorias", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, mockTx, mockHasher, service := setupUserFullService()

		invalidUser := &model_user_full.UserFull{
			User:    &model_user.User{Username: "user", Email: "user@test.com"},
			Address: &model_address.Address{Street: "Rua Teste"},
			Contact: &model_contact.Contact{Phone: "1112345678"},
		}

		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		result, err := service.CreateFull(ctx, invalidUser)

		assert.Nil(t, result)
		assert.ErrorContains(t, err, "pelo menos uma categoria")
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
	})

	t.Run("transacao_nil", func(t *testing.T) {
		mockUserRepo, _, _, _, _, _, service := setupUserFullService()

		userFull := &model_user_full.UserFull{
			User:    &model_user.User{Username: "user", Email: "user@test.com"},
			Address: &model_address.Address{Street: "Rua Teste"},
			Contact: &model_contact.Contact{Phone: "1112345678"},
			Categories: []model_user_cat_rel.UserCategory{
				{ID: 1},
			},
		}

		mockUserRepo.On("BeginTx", mock.Anything).Return(nil, nil)

		_, err := service.CreateFull(ctx, userFull)

		assert.Error(t, err)
		assert.EqualError(t, err, "transação inválida")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("erro_ao_iniciar_transacao", func(t *testing.T) {
		mockUserRepo, _, _, _, _, _, service := setupUserFullService()

		userFull := &model_user_full.UserFull{
			User:    &model_user.User{Username: "user", Email: "user@test.com"},
			Address: &model_address.Address{Street: "Rua Teste"},
			Contact: &model_contact.Contact{Phone: "1112345678"},
			Categories: []model_user_cat_rel.UserCategory{
				{ID: 1},
			},
		}

		mockUserRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("falha na transação"))

		_, err := service.CreateFull(ctx, userFull)

		assert.ErrorContains(t, err, "erro ao iniciar transação")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("sucesso_na_criacao_completa", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, service := setupUserFullService()

		userFull := &model_user_full.UserFull{
			User:    &model_user.User{Username: "user", Email: "user@test.com", Password: "senha123"},
			Address: &model_address.Address{Street: "Rua Teste"},
			Contact: &model_contact.Contact{Phone: "1112345678"},
			Categories: []model_user_cat_rel.UserCategory{
				{ID: 1},
			},
		}

		hashedPassword := "hashed_senha"
		mockHasher.On("Hash", "senha123").Return(hashedPassword, nil)

		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: "1", Username: "user", Password: hashedPassword}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil)

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: "1", CategoryID: 1}, nil)

		mockTx.On("Commit", mock.Anything).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.NoError(t, err)
		assert.Equal(t, "1", result.User.UID)
		assert.Equal(t, hashedPassword, result.User.Password)

		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}
