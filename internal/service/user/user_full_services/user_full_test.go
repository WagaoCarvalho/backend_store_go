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
	model_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	model_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	model_user_full "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(_, _ string) error {
	return nil
}

func createValidUserFull() *model_user_full.UserFull {
	return &model_user_full.UserFull{
		User: &model_user.User{
			Username: "usuario_valido",
			Email:    "email@valido.com",
			Password: "Senha123!",
		},
		Address: &model_address.Address{
			Street:     "Rua Valida",
			City:       "Cidade Valida",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		},
		Contact: &model_contact.Contact{
			ContactName: "Contato Valido",
			Phone:       "1112345678",
			Email:       "contato@valido.com",
		},
		Categories: []model_user_categories.UserCategory{
			{ID: 1},
		},
	}
}

func TestUserService_CreateFull(t *testing.T) {
	setup := func() (
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

		userService := NewUserFullService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService
	}

	t.Run("transação_nil", func(t *testing.T) {
		mockUserRepo, _, _, _, _, mockHasher, userService := setup()

		mockHasher.On("Hash", "Senha123").Return("senha-hash", nil)

		mockUserRepo.On("BeginTx", mock.Anything).Return(nil, nil)

		userFull := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "usuario_teste",
				Email:    "email@invalido.com",
				Password: "Senha123",
				Status:   true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
				{ID: 2},
			},
		}

		_, err := userService.CreateFull(context.Background(), userFull)

		assert.Error(t, err)
		assert.EqualError(t, err, "transação inválida")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("userFull_ou_user_nil", func(t *testing.T) {
		_, _, _, _, _, _, userService := setup()

		// Caso userFull == nil
		result, err := userService.CreateFull(context.Background(), nil)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "userFull é nulo")

		// Caso userFull.User == nil
		invalidUser := &model_user_full.UserFull{}
		result, err = userService.CreateFull(context.Background(), invalidUser)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "usuário é obrigatório")
	})

	t.Run("erro_ao_hash_senha", func(t *testing.T) {
		_, _, _, _, _, mockHasher, userService := setup()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "usuario_teste",
				Email:    "email@invalido.com",
				Password: "Senha123",
				Status:   true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
				{ID: 2},
			},
		}

		hashErr := errors.New("falha ao hashear")

		mockHasher.On("Hash", "Senha123").Return("", hashErr).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")

		mockHasher.AssertExpectations(t)
	})

	t.Run("erro ao iniciar transação", func(t *testing.T) {
		mockUserRepo, _, _, _, _, mockHasher, userService := setup()

		userFull := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "usuario_teste",
				Email:    "email@invalido.com",
				Password: "Senha123",
				Status:   true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
				{ID: 2},
			},
		}

		mockHasher.On("Hash", "Senha123").Return("hashedSenha123", nil).Once()

		mockUserRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("falha na transação")).Once()

		_, err := userService.CreateFull(context.Background(), userFull)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao iniciar transação")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("erro_ao_fazer_rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "Senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil)

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar relação"))

		mockTx.On("Rollback", mock.Anything).Return(errors.New("erro ao dar rollback"))

		userFull := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "usuario_teste",
				Email:    "usuario@teste.com",
				Password: "Senha123",
				Status:   true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "1112345678",
				Email:       "joao@teste.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
		}

		_, err := userService.CreateFull(context.Background(), userFull)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação")
		assert.Contains(t, err.Error(), "rollback error")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_commitar_transacao", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "Senha123").Return("senha-hash", nil)

		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{
				UID:      1,
				Username: "teste",
				Email:    "teste@example.com",
				Password: "senha-hash",
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil)

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil)

		mockTx.On("Commit", mock.Anything).Return(errors.New("erro ao commitar transação"))

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		userFull := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "teste",
				Email:    "teste@example.com",
				Password: "Senha123",
			},
			Address: &model_address.Address{
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
		}

		_, err := userService.CreateFull(context.Background(), userFull)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_fazer_rollback_apos_commit_falhar", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		userFull := createValidUserFull()

		mockHasher.On("Hash", mock.Anything).Return("hashed_password", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).Return(&model_user.User{UID: 1}, nil)
		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).Return(&model_address.Address{ID: 1}, nil)
		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).Return(&model_contact.Contact{ID: 1}, nil)
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).Return(nil, nil)

		commitError := errors.New("erro no commit")
		mockTx.On("Commit", mock.Anything).Return(commitError)

		rollbackError := errors.New("erro no rollback")
		mockTx.On("Rollback", mock.Anything).Return(rollbackError)

		_, err := userService.CreateFull(context.Background(), userFull)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")
		assert.Contains(t, err.Error(), "erro no commit")
		assert.Contains(t, err.Error(), "rollback error")
		assert.Contains(t, err.Error(), "erro no rollback")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_endereco", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, _, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "Senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{
				UID:      1,
				Username: "teste",
				Email:    "teste@example.com",
				Password: "senha-hash",
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar endereço"))

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		userFull := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "teste",
				Email:    "teste@example.com",
				Password: "Senha123",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
			Contact: &model_contact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
		}

		_, err := userService.CreateFull(context.Background(), userFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_contato", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "Senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{
				UID:      1,
				Username: "teste",
				Email:    "teste@example.com",
				Password: "senha-hash",
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{
				ID:         1,
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar contato"))

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		userFull := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "teste",
				Email:    "teste@example.com",
				Password: "Senha123",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
		}

		_, err := userService.CreateFull(context.Background(), userFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_usuario_na_transacao", func(t *testing.T) {
		mockUserRepo, _, _, _, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "Senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(user *model_user.User) bool {
			return user.Email == "test@example.com" && user.Password == "senha-hash"
		})).Return(nil, errors.New("falha ao criar usuário")).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Username: "vsdvvfvf",
				Email:    "test@example.com",
				Password: "Senha123",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
		}

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "falha ao criar usuário")

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("sucesso_na_criacao_completa", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Email:    "test@example.com",
				Password: "Senha123",
				Username: "teste",
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
		}

		expectedUser := &model_user.User{
			UID:      1,
			Email:    "test@example.com",
			Username: "teste",
		}

		mockHasher.On("Hash", "Senha123").
			Return("hashed", nil).Once()

		mockUserRepo.On("BeginTx", mock.Anything).
			Return(mockTx, nil).Once()

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(u *model_user.User) bool {
			return u.Email == "test@example.com" && u.Password == "hashed"
		})).Return(expectedUser, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(addr *model_address.Address) bool {
			return addr.City == "Cidade B" && addr.PostalCode == "12345678"
		})).Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.Phone == "1234567895" && c.ContactName == "Ari"
		})).Return(&model_contact.Contact{ID: 1}, nil).Once()

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil).Once()

		mockTx.On("Commit", mock.Anything).Return(nil).Once()

		result, err := userService.CreateFull(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.UID, result.User.UID)
		assert.Equal(t, expectedUser.Email, result.User.Email)
		assert.Equal(t, expectedUser.Username, result.User.Username)

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
	})

	t.Run("falha_validacao_do_endereco_faz_rollback", func(t *testing.T) {
		mockUserRepo, _, _, _, mockTx, mockHasher, userService := setup()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Email:    "test@example.com",
				Password: "Senha123",
				Username: "teste",
			},
			Address: &model_address.Address{
				Street: "", // força falha de validação
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
		}

		mockHasher.On("Hash", "Senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "endereço inválido")

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("falha_validacao_do_contato_faz_rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, _, _, mockTx, mockHasher, userService := setup()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Email:    "test@example.com",
				Password: "Senha123",
				Username: "teste",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				Phone: "invalido", // força erro de validação
			},
		}

		mockHasher.On("Hash", "Senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato inválido")

		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("falha_validacao_relacao_usuario_categoria_faz_rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, mockTx, mockHasher, userService := setup()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Email:    "test@example.com",
				Password: "Senha123",
				Username: "teste",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 0},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João",
				Phone:       "1234567890",
				Email:       "joao@example.com",
			},
		}

		mockHasher.On("Hash", "Senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "relação usuário-categoria inválida")

		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("panic_faz_rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, mockTx, mockHasher, userService := setup()

		user := &model_user_full.UserFull{
			User: &model_user.User{
				Email:    "test@example.com",
				Password: "Senha123",
				Username: "teste",
			},
			Categories: []model_user_categories.UserCategory{
				{ID: 1},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
		}

		mockHasher.On("Hash", "Senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Run(func(_ mock.Arguments) {
				panic("panic simulado")
			}).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		assert.Panics(t, func() {
			_, _ = userService.CreateFull(context.Background(), user)
		})

		mockTx.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}
