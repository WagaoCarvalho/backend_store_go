package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	modelsAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelsContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	modelUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	modelsUserCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	modelsUserContactRel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relations"
	modelsUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"

	mockAuth "github.com/WagaoCarvalho/backend_store_go/infra/mock/auth"
	mockTX "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo"
	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/address"
	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	mockUserCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
	mockUserContactRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/user"
)

func TestUserFullService_CreateFull(t *testing.T) {
	ctx := context.Background()

	// ------------------------
	// Grupo: User
	// ------------------------
	t.Run("User: falha quando userFull é nil", func(t *testing.T) {
		service := NewUserFullService(
			new(mockUser.MockUserFullRepository),
			nil, nil, nil, nil,
			new(mockAuth.MockHasher),
		)

		result, err := service.CreateFull(ctx, nil)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrInvalidData.Error())
	})

	t.Run("User: falha ao validar usuário inválido", func(t *testing.T) {
		service := NewUserFullService(
			new(mockUser.MockUserFullRepository),
			nil, nil, nil, nil,
			new(mockAuth.MockHasher),
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{Password: "123"},
		}

		result, err := service.CreateFull(ctx, userFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrInvalidData.Error())
	})

	t.Run("User: falha ao hashear senha", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockHasher := new(mockAuth.MockHasher)

		service := NewUserFullService(
			mockRepoUser,
			nil, nil, nil, nil,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockHasher.On("Hash", "SenhaValida@123").Return("", errors.New("hash error"))

		result, err := service.CreateFull(ctx, userFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")
		mockHasher.AssertExpectations(t)
	})

	t.Run("User: falha ao criar usuário no banco", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		// fluxo: BeginTx -> Hash -> erro no Create User
		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(nil, errors.New("erro ao criar usuário"))
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar usuário")
	})

	// ------------------------
	// Grupo: Address
	// ------------------------
	t.Run("Address: falha ao validar endereço inválido", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			nil, nil, nil,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		// Mock necessário para evitar panic
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "endereço inválido")
	})

	t.Run("Address: falha ao criar endereço no banco", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(nil, errors.New("erro ao criar endereço"))
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")
	})

	// ------------------------
	// Grupo: Contact
	// ------------------------
	t.Run("Contact: falha ao validar contato inválido", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			nil, nil,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "",
				Email:       "csacc@gmail.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		tx.On("Rollback", ctx).Return(nil)

		// Se a validação permitir continuar, mock CreateTx do contato
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(nil, errors.New("não deveria chegar aqui")).Maybe()

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato inválido")
	})

	t.Run("Contact: falha ao criar contato no banco", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(nil, errors.New("erro ao criar contato"))
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")
	})

	// ------------------------
	// Grupo: Relations (User-Contact e User-Category)
	// ------------------------
	t.Run("Relations: sucesso ao criar relações user-contact e user-category", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "Sasss@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "Sasss@123").Return("hashed", nil)
		createdUser := *userFull.User
		createdUser.Password = "hashed"
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(&createdUser, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		mockRepoUserContactRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations")).
			Return(&modelsUserContactRel.UserContactRelations{
				UserID:    createdUser.UID,
				ContactID: userFull.Contact.ID,
				CreatedAt: time.Now(),
			}, nil)
		mockRepoUserCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserCategoryRelations")).
			Return(&modelsUserCatRel.UserCategoryRelations{
				UserID:     createdUser.UID,
				CategoryID: 100,
				CreatedAt:  time.Now(),
			}, nil)

		tx.On("Commit", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "hashed", result.User.Password)
	})

	t.Run("UserContactRelation: falha ao validar relação inválida", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			nil,
			mockRepoUserContactRel,
			mockHasher,
		)

		// Dados inválidos para a relação
		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		// Mock básico para criar usuário e outros dados necessários antes da relação
		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		tx.On("Rollback", ctx).Return(nil)

		// Força relação inválida: por exemplo, IDs zerados
		userFull.Contact.ID = 0
		userFull.User.UID = 0

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "relação usuário-contato inválida")
	})

	t.Run("UserContactRelation: falha ao criar relação no repo", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			nil,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		// Configurações de mocks para criar usuário, endereço e contato
		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		tx.On("Rollback", ctx).Return(nil)

		// Mock do CreateTx do repositório de relação para retornar erro
		mockRepoUserContactRel.On(
			"CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations"),
		).Return(nil, errors.New("erro ao criar relação"))

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação")
	})

	t.Run("UserCategoryRelation: falha ao validar relação categoria inválida", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 0}, // ID inválido
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		mockRepoUserContactRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations")).
			Return(&modelsUserContactRel.UserContactRelations{
				UserID:    userFull.User.UID,
				ContactID: userFull.Contact.ID,
				CreatedAt: time.Now(),
			}, nil)
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "relação usuário-categoria inválida")
	})

	t.Run("UserCategoryRelation: falha ao criar relação categoria no repo", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		mockRepoUserContactRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations")).
			Return(&modelsUserContactRel.UserContactRelations{
				UserID:    userFull.User.UID,
				ContactID: userFull.Contact.ID,
				CreatedAt: time.Now(),
			}, nil)
		mockRepoUserCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserCategoryRelations")).
			Return(nil, errors.New("erro ao criar relação categoria"))
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação categoria")
	})

	// ------------------------
	// Grupo: Transaction
	// ------------------------

	t.Run("Transaction: falha ao iniciar transação", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockHasher := new(mockAuth.MockHasher)

		service := NewUserFullService(
			mockRepoUser,
			new(mockAddress.MockAddressRepository),
			new(mockContact.MockContactRepository),
			new(mockUserCatRel.MockUserCategoryRelationRepo),
			new(mockUserContactRel.MockUserContactRelationRepo),
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("BeginTx", ctx).Return(nil, errors.New("begin error"))

		result, err := service.CreateFull(ctx, userFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao iniciar transação")
	})

	t.Run("Transaction: falha quando transação é nil", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockHasher := new(mockAuth.MockHasher)

		service := NewUserFullService(
			mockRepoUser,
			new(mockAddress.MockAddressRepository),
			new(mockContact.MockContactRepository),
			new(mockUserCatRel.MockUserCategoryRelationRepo),
			new(mockUserContactRel.MockUserContactRelationRepo),
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("BeginTx", ctx).Return(nil, nil)

		result, err := service.CreateFull(ctx, userFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transação inválida")
	})

	t.Run("Transaction: falha no rollback quando há erro", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			new(mockAddress.MockAddressRepository),
			new(mockContact.MockContactRepository),
			new(mockUserCatRel.MockUserCategoryRelationRepo),
			new(mockUserContactRel.MockUserContactRelationRepo),
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(nil, errors.New("erro ao criar usuário"))
		tx.On("Rollback", ctx).Return(errors.New("erro no rollback"))

		result, err := service.CreateFull(ctx, userFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar usuário")
		assert.Contains(t, err.Error(), "rollback error")
	})

	t.Run("Transaction: faz rollback em caso de panic", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			new(mockContact.MockContactRepository),
			new(mockUserCatRel.MockUserCategoryRelationRepo),
			new(mockUserContactRel.MockUserContactRelationRepo),
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)

		// Simula panic durante a criação do usuário
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Run(func(args mock.Arguments) {
			panic("panic simulado durante criação do usuário")
		}).Return(nil, nil)

		// Rollback deve ser chamado pelo defer
		tx.On("Rollback", ctx).Return(nil)

		defer func() {
			if r := recover(); r != nil {
				// Verifica se rollback foi chamado
				tx.AssertCalled(t, "Rollback", ctx)
				assert.Equal(t, "panic simulado durante criação do usuário", r)
			} else {
				t.Errorf("Esperado panic, mas não ocorreu")
			}
		}()

		// Esse call vai disparar o panic
		_, _ = service.CreateFull(ctx, userFull)

		t.Errorf("Não deveria chegar aqui, panic esperado")
	})

	t.Run("Transaction: falha no commit e rollback também falha", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "SenhaValida@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "SenhaValida@123").Return("hashed", nil)

		createdUser := *userFull.User
		createdUser.Password = "hashed"
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(&createdUser, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		mockRepoUserContactRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations")).
			Return(&modelsUserContactRel.UserContactRelations{
				UserID:    createdUser.UID,
				ContactID: userFull.Contact.ID,
				CreatedAt: time.Now(),
			}, nil)
		mockRepoUserCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserCategoryRelations")).
			Return(&modelsUserCatRel.UserCategoryRelations{
				UserID:     createdUser.UID,
				CategoryID: 100,
				CreatedAt:  time.Now(),
			}, nil)

		// Cenário: commit falha e rollback também falha
		commitErr := errors.New("erro no commit")
		rollbackErr := errors.New("erro no rollback")
		tx.On("Commit", ctx).Return(commitErr)
		tx.On("Rollback", ctx).Return(rollbackErr)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")
		assert.Contains(t, err.Error(), "rollback error")
		assert.Contains(t, err.Error(), commitErr.Error())
		assert.Contains(t, err.Error(), rollbackErr.Error())
	})

	t.Run("Transaction: falha no commit e rollback com sucesso", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "Senha123!",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "Senha123!").Return("hashed", nil)
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(userFull.User, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		mockRepoUserContactRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations")).
			Return(&modelsUserContactRel.UserContactRelations{
				UserID:    userFull.User.UID,
				ContactID: userFull.Contact.ID,
				CreatedAt: time.Now(),
			}, nil)
		mockRepoUserCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserCategoryRelations")).
			Return(&modelsUserCatRel.UserCategoryRelations{
				UserID:     userFull.User.UID,
				CategoryID: 100,
				CreatedAt:  time.Now(),
			}, nil)

		commitErr := errors.New("erro no commit")
		tx.On("Commit", ctx).Return(commitErr)
		// Rollback deve ser chamado, mas sem erro
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")
		assert.Contains(t, err.Error(), commitErr.Error())

		// Valida se Commit e Rollback foram chamados
		tx.AssertCalled(t, "Commit", ctx)
		tx.AssertCalled(t, "Rollback", ctx)
	})

	// ------------------------
	// Grupo: Success
	// ------------------------
	t.Run("Success: cria userFull com múltiplas categorias", func(t *testing.T) {
		mockRepoUser := new(mockUser.MockUserFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoUserCatRel := new(mockUserCatRel.MockUserCategoryRelationRepo)
		mockRepoUserContactRel := new(mockUserContactRel.MockUserContactRelationRepo)
		mockHasher := new(mockAuth.MockHasher)
		tx := new(mockTX.MockTx)

		service := NewUserFullService(
			mockRepoUser,
			mockRepoAddress,
			mockRepoContact,
			mockRepoUserCatRel,
			mockRepoUserContactRel,
			mockHasher,
		)

		userFull := &modelsUserFull.UserFull{
			User: &modelsUser.User{
				UID:      1,
				Username: "Walla",
				Status:   true,
				Email:    "test@example.com",
				Password: "Sasss@123",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				UserID:       utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelUserCategories.UserCategory{
				{ID: 100},
				{ID: 200},
				{ID: 300},
			},
		}

		mockRepoUser.On("BeginTx", ctx).Return(tx, nil)
		mockHasher.On("Hash", "Sasss@123").Return("hashed", nil)
		createdUser := *userFull.User
		createdUser.Password = "hashed"
		mockRepoUser.On("CreateTx", ctx, tx, userFull.User).Return(&createdUser, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, userFull.Address).Return(userFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, userFull.Contact).Return(userFull.Contact, nil)
		mockRepoUserContactRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserContactRelations")).
			Return(&modelsUserContactRel.UserContactRelations{
				UserID:    createdUser.UID,
				ContactID: userFull.Contact.ID,
				CreatedAt: time.Now(),
			}, nil)

		// Mock para múltiplas categorias
		for i := 0; i < len(userFull.Categories); i++ {
			mockRepoUserCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.UserCategoryRelations")).
				Return(&modelsUserCatRel.UserCategoryRelations{
					UserID:     createdUser.UID,
					CategoryID: int64(userFull.Categories[i].ID),
					CreatedAt:  time.Now(),
				}, nil)
		}

		tx.On("Commit", ctx).Return(nil)

		result, err := service.CreateFull(ctx, userFull)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "hashed", result.User.Password)
		assert.Len(t, result.Categories, 3)
	})

}
