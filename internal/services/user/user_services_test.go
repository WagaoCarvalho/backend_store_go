package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	model_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	model_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	model_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	model_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repo_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	repo_contact "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repo_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	repo_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	"github.com/sirupsen/logrus"
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
	// Implementado apenas para satisfazer a interface
	return nil
}

func TestUserService_Create(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // logger real

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("erro ao hashear senha", func(t *testing.T) {
		mockUserRepo, _, _, _, mockHasher, userService := setup()

		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senhaInvalidaParaHash",
		}

		mockHasher.On("Hash", "senhaInvalidaParaHash").Return("", errors.New("falha no hash")).Once()

		_, err := userService.Create(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")
		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("sucesso ao criar usuário com todos os dados", func(t *testing.T) {
		mockUserRepo, _, _, _, mockHasher, userService := setup()

		newUser := &model_user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "senha123",
			Status:   true,
		}

		hashed := "hashedSenha123"
		mockHasher.On("Hash", "senha123").Return(hashed, nil).Once()

		createdUser := &model_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashed,
			Status:   true,
		}

		mockUserRepo.
			On("Create", mock.Anything, mock.MatchedBy(func(u *model_user.User) bool {
				return u.Email == newUser.Email && u.Password == hashed
			})).
			Run(func(args mock.Arguments) {
				args.Get(1).(*model_user.User).UID = 1
			}).
			Return(createdUser, nil)

		result, err := userService.Create(context.Background(), newUser)

		assert.NoError(t, err)
		assert.Equal(t, createdUser, result)
		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
	t.Run("erro ao criar usuário", func(t *testing.T) {
		mockUserRepo, _, _, _, _, userService := setup()

		newUser := model_user.User{Email: "test@example.com"}
		mockUserRepo.On("Create", mock.Anything, &newUser).Return(nil, errors.New("erro no banco de dados"))

		_, err := userService.Create(context.Background(), &newUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar usuário")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("usuário criado é nulo", func(t *testing.T) {
		mockUserRepo, _, _, _, _, userService := setup()

		user := &model_user.User{
			Email: "valid@email.com",
		}

		mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil, nil)

		_, err := userService.Create(context.Background(), user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "usuário criado é nulo")
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("email inválido", func(t *testing.T) {
		_, _, _, _, _, userService := setup()

		_, err := userService.Create(context.Background(), &model_user.User{Email: "email-invalido"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
	})

}

func TestUserService_CreateFull(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*repo_user.MockTx,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)
		mockTx := new(repo_user.MockTx)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService
	}

	t.Run("email inválido", func(t *testing.T) {
		_, _, _, _, _, _, userService := setup()
		user := &model_user.User{Email: "email-invalido"}

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidEmail, err)
	})

	t.Run("transação_nil", func(t *testing.T) {
		mockUserRepo, _, _, _, _, mockHasher, userService := setup()

		// Configura hasher para retornar senha hash válida
		mockHasher.On("Hash", "senha123").Return("senha-hash", nil)

		// Simula retorno de transação nil sem erro
		mockUserRepo.On("BeginTx", mock.Anything).Return(nil, nil)

		user := &model_user.User{
			Username: "teste",
			Email:    "teste@example.com",
			Password: "senha123",
		}

		_, err := userService.CreateFull(context.Background(), user)

		assert.Error(t, err)
		assert.EqualError(t, err, "transação inválida")
		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("erro ao hashear senha", func(t *testing.T) {
		mockUserRepo, _, _, _, _, mockHasher, userService := setup()
		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senha123",
		}

		mockHasher.On("Hash", "senha123").Return("", errors.New("falha no hash")).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao hashear senha")

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("erro ao iniciar transação", func(t *testing.T) {
		mockUserRepo, _, _, _, _, mockHasher, userService := setup()
		user := &model_user.User{Email: "test@example.com"}

		mockHasher.On("Hash", "").Return("", nil).Maybe()
		mockUserRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("falha na transação")).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao iniciar transação")

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("erro_ao_fazer_rollback", func(t *testing.T) {
		mockUserRepo, _, _, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar relação"))

		mockUserRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
			Return(&model_user.User{UID: 1}, nil)

		mockTx.On("Rollback", mock.Anything).Return(errors.New("erro ao dar rollback"))

		user := &model_user.User{
			Username: "teste",
			Email:    "teste@example.com",
			Password: "senha123",
			Categories: []model_categories.UserCategory{
				{ID: 1},
			},
		}

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rollback error")
		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_commitar_transacao", func(t *testing.T) {
		mockUserRepo, _, _, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{
				UID:      1,
				Username: "teste",
				Email:    "teste@example.com",
				Password: "senha-hash",
			}, nil)
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil)

		// Mock para Commit com erro
		mockTx.On("Commit", mock.Anything).Return(errors.New("erro ao commitar transação"))

		user := &model_user.User{
			Username: "teste",
			Email:    "teste@example.com",
			Password: "senha123",
			Categories: []model_categories.UserCategory{
				{ID: 1},
			},
		}

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_endereco", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, _, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		// Mock para CreateTx do UserRepository
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{
				UID:      1,
				Username: "teste",
				Email:    "teste@example.com",
				Password: "senha-hash",
			}, nil)

		// Mock para CreateTx do AddressRepository que simula erro
		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar endereço"))

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.AnythingOfType("*models.UserCategoryRelations")).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		user := &model_user.User{
			Username: "teste",
			Email:    "teste@example.com",
			Password: "senha123",
			Categories: []model_categories.UserCategory{
				{ID: 1},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "99999999",
			},
		}

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_contato", func(t *testing.T) {
		mockUserRepo, _, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		mockHasher.On("Hash", "senha123").Return("senha-hash", nil)
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{
				UID:      1,
				Username: "teste",
				Email:    "teste@example.com",
				Password: "senha-hash",
			}, nil)

		// Mock erro ao criar contato
		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar contato"))

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.AnythingOfType("*models.UserCategoryRelations")).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		user := &model_user.User{
			Username: "teste",
			Email:    "teste@example.com",
			Password: "senha123",
			Categories: []model_categories.UserCategory{
				{ID: 1},
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
		}

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")

		mockHasher.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro ao criar usuário na transação", func(t *testing.T) {
		mockUserRepo, _, _, _, mockTx, mockHasher, userService := setup()
		user := &model_user.User{Email: "test@example.com"}

		mockHasher.On("Hash", "").Return("", nil).Maybe()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, user).
			Return(nil, errors.New("falha ao criar usuário")).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()
		mockTx.On("Commit", mock.Anything).Return(nil).Maybe()

		_, err := userService.CreateFull(context.Background(), user)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "falha ao criar usuário")

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("sucesso na criação completa", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, mockHasher, userService := setup()

		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senha123",
			Username: "teste",
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "99999999",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []model_categories.UserCategory{{ID: 1}},
		}

		expectedUser := &model_user.User{
			UID:      1,
			Email:    "test@example.com",
			Username: "teste",
		}

		mockHasher.On("Hash", "senha123").
			Return("hashed", nil).
			Once()

		mockUserRepo.On("BeginTx", mock.Anything).
			Return(mockTx, nil).
			Once()

		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(expectedUser, nil).
			Once()

		mockAddressRepo.
			On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(addr *model_address.Address) bool {
				return addr.City == "Cidade B" && addr.PostalCode == "99999999"
			})).
			Return(&model_address.Address{ID: 1}, nil).
			Once()

		mockContactRepo.
			On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(c *model_contact.Contact) bool {
				return c.Phone == "1234567895" && c.ContactName == "Ari"
			})).
			Return(&model_contact.Contact{ID: 1}, nil).
			Once()

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user_cat_rel.UserCategoryRelations{UserID: 1, CategoryID: 1}, nil).
			Once()

		mockTx.On("Commit", mock.Anything).Return(nil).Once()

		result, err := userService.CreateFull(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser.UID, result.UID)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Username, result.Username)

		mockUserRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
	})

	t.Run("falha validação do endereço faz rollback", func(t *testing.T) {
		mockUserRepo, _, _, _, mockTx, mockHasher, userService := setup()

		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senha123",
			Username: "teste",
			Address: &model_address.Address{
				Street: "", // força falha de validação
			},
		}

		mockHasher.On("Hash", "senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.ErrorContains(t, err, "endereço inválido")

		mockTx.AssertExpectations(t)
	})

	t.Run("falha validação do contato faz rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, _, _, mockTx, mockHasher, userService := setup()

		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senha123",
			Username: "teste",
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "99999999",
			},
			Contact: &model_contact.Contact{
				Phone: "invalido", // força erro de validação
			},
			Categories: []model_categories.UserCategory{{ID: 1}},
		}

		mockHasher.On("Hash", "senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.ErrorContains(t, err, "contato inválido")

		mockTx.AssertExpectations(t)
	})

	t.Run("falha validação da relação usuário-categoria faz rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, mockTx, mockHasher, userService := setup()

		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senha123",
			Username: "teste",
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "99999999",
			},
			Contact: &model_contact.Contact{
				ContactName: "João",
				Phone:       "1234567890",
				Email:       "joao@example.com",
			},
			Categories: []model_categories.UserCategory{
				{ID: 0}, // ID inválido → força falha de validação
			},
		}

		mockHasher.On("Hash", "senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := userService.CreateFull(context.Background(), user)
		assert.ErrorContains(t, err, "relação usuário-categoria inválida")

		mockTx.AssertExpectations(t)
	})

	t.Run("panic faz rollback", func(t *testing.T) {
		mockUserRepo, mockAddressRepo, mockContactRepo, _, mockTx, mockHasher, userService := setup()

		user := &model_user.User{
			Email:    "test@example.com",
			Password: "senha123",
			Username: "teste",
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "99999999",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []model_categories.UserCategory{
				{ID: 1},
			},
		}

		mockHasher.On("Hash", "senha123").Return("hashed", nil).Once()
		mockUserRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockUserRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_user.User{UID: 1, Email: "test@example.com", Username: "teste"}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Run(func(args mock.Arguments) {
				panic("panic simulado")
			}).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		assert.Panics(t, func() {
			_, _ = userService.CreateFull(context.Background(), user)
		})

		mockTx.AssertExpectations(t)
	})

}

func TestUserService_GetAll(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve retornar todos os usuários com sucesso", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		expectedUsers := []*model_user.User{
			{UID: 1, Username: "user1", Email: "user1@example.com", Status: true},
			{UID: 2, Username: "user2", Email: "user2@example.com", Status: false},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedUsers, nil)

		users, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao falhar no repositório", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetAll", mock.Anything).Return(nil, fmt.Errorf("erro ao acessar o banco"))

		users, err := service.GetAll(context.Background())

		assert.ErrorContains(t, err, "erro ao acessar o banco")
		assert.Nil(t, users)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve retornar usuário quando encontrado", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		expectedUser := &model_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedUser, nil)

		user, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando usuário não existe", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, fmt.Errorf("usuário não encontrado"))

		user, err := service.GetByID(context.Background(), 999)

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetVersionByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve retornar versão quando usuário for encontrado", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil)

		version, err := service.GetVersionByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(999)).Return(
			int64(0),
			repo_user.ErrUserNotFound,
		)

		version, err := service.GetVersionByID(context.Background(), 999)

		assert.ErrorIs(t, err, repo_user.ErrUserNotFound)
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico quando falhar no repositório", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(2)).Return(
			int64(0),
			fmt.Errorf("erro no banco de dados"),
		)

		version, err := service.GetVersionByID(context.Background(), 2)

		assert.ErrorContains(t, err, "versão inválida")
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve retornar usuário quando encontrado por e-mail", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		expectedUser := &model_user.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}

		mockRepo.On("GetByEmail", mock.Anything, "user1@example.com").Return(expectedUser, nil)

		user, err := service.GetByEmail(context.Background(), "user1@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando repositório falha ao buscar por e-mail", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetByEmail", mock.Anything, "inexistente@example.com").Return(
			nil,
			fmt.Errorf("usuário não encontrado"),
		)

		user, err := service.GetByEmail(context.Background(), "inexistente@example.com")

		assert.ErrorContains(t, err, "usuário não encontrado")
		assert.Nil(t, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve retornar erro ao atualizar com e-mail inválido", func(t *testing.T) {
		_, _, _, _, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "email_invalido",
			Version: 1,
		}

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, ErrInvalidEmail)
		assert.Nil(t, updatedUser)
	})

	t.Run("Deve retornar erro ao atualizar com versão inválida", func(t *testing.T) {
		_, _, _, _, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 0,
		}

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, ErrInvalidVersion)
		assert.Nil(t, updatedUser)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, repo_user.ErrUserNotFound)

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, repo_user.ErrUserNotFound)
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 2,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, repo_user.ErrVersionConflict)

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorIs(t, err, repo_user.ErrVersionConflict)
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao atualizar", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		user := &model_user.User{
			UID:     1,
			Email:   "user@example.com",
			Version: 1,
		}

		mockRepo.On("Update", mock.Anything, user).Return(nil, fmt.Errorf("erro interno"))

		updatedUser, err := service.Update(context.Background(), user)

		assert.ErrorContains(t, err, "erro ao atualizar usuário")
		assert.Nil(t, updatedUser)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve atualizar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		user := &model_user.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  1,
		}

		expected := &model_user.User{
			UID:      1,
			Username: "usuario",
			Email:    "user@example.com",
			Version:  2,
		}

		mockRepo.On("Update", mock.Anything, user).Return(expected, nil)

		updatedUser, err := service.Update(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, expected, updatedUser)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Disable(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve desativar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		original := &model_user.User{
			UID:     1,
			Email:   "test@example.com",
			Status:  true,
			Version: 1,
		}
		updated := *original
		updated.Status = false

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(original, nil).Once()
		mockRepo.On("Update", mock.Anything, &updated).Return(&updated, nil).Once()

		err := service.Disable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao buscar usuário", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return(nil, fmt.Errorf("erro ao buscar")).Once()

		err := service.Disable(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao buscar usuário")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao atualizar status para falso", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		original := &model_user.User{
			UID:     3,
			Email:   "fail@example.com",
			Status:  true,
			Version: 2,
		}
		updated := *original
		updated.Status = false

		mockRepo.On("GetByID", mock.Anything, int64(3)).Return(original, nil).Once()
		mockRepo.On("Update", mock.Anything, &updated).
			Return(nil, fmt.Errorf("falha ao atualizar")).Once()

		err := service.Disable(context.Background(), 3)

		assert.ErrorContains(t, err, "erro ao desabilitar usuário")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Enable(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve ativar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		original := &model_user.User{
			UID:     1,
			Email:   "test@example.com",
			Status:  false,
			Version: 1,
		}
		updated := *original
		updated.Status = true

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(original, nil).Once()
		mockRepo.On("Update", mock.Anything, &updated).Return(&updated, nil).Once()

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao buscar usuário", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return(nil, fmt.Errorf("erro ao buscar")).Once()

		err := service.Enable(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao buscar usuário")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao atualizar status para true", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		original := &model_user.User{
			UID:     3,
			Email:   "fail@example.com",
			Status:  false,
			Version: 2,
		}
		updated := *original
		updated.Status = true

		mockRepo.On("GetByID", mock.Anything, int64(3)).Return(original, nil).Once()
		mockRepo.On("Update", mock.Anything, &updated).
			Return(nil, fmt.Errorf("falha ao atualizar")).Once()

		err := service.Enable(context.Background(), 3)

		assert.ErrorContains(t, err, "erro ao habilitar usuário")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_user.MockUserRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_user_cat_rel.MockUserCategoryRelationRepo,
		*MockHasher,
		UserService,
	) {
		mockUserRepo := new(repo_user.MockUserRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_user_cat_rel.MockUserCategoryRelationRepo)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
			mockHasher,
		)

		return mockUserRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockHasher, userService
	}

	t.Run("Deve deletar usuário com sucesso", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao falhar na deleção", func(t *testing.T) {
		mockRepo, _, _, _, _, service := setup()

		mockRepo.On("Delete", mock.Anything, int64(2)).Return(fmt.Errorf("erro no banco"))

		err := service.Delete(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao deletar usuário")
		mockRepo.AssertExpectations(t)
	})
}
