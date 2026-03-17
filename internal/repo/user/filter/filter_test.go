package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	filterUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUser_Filter(t *testing.T) {
	t.Run("successfully get all users", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Preparar valores
		description := "Descrição do usuário"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // id/UID
			mock.AnythingOfType("*string"),    // username
			mock.AnythingOfType("*string"),    // email
			mock.AnythingOfType("*string"),    // password_hash
			mock.AnythingOfType("*string"),    // description
			mock.AnythingOfType("*bool"),      // status
			mock.AnythingOfType("*int"),       // version
			mock.AnythingOfType("*time.Time"), // created_at
			mock.AnythingOfType("*time.Time"), // updated_at
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "john_doe"
			*args[2].(*string) = "john@example.com"
			*args[3].(*string) = "hashed_password_123"
			*args[4].(*string) = description
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now
			*args[8].(*time.Time) = now
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(1), result[0].UID)
		assert.Equal(t, "john_doe", result[0].Username)
		assert.Equal(t, "john@example.com", result[0].Email)
		assert.Equal(t, "hashed_password_123", result[0].Password)
		assert.Equal(t, description, result[0].Description)
		assert.True(t, result[0].Status)
		assert.Equal(t, 1, result[0].Version)
		assert.WithinDuration(t, now, result[0].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[0].UpdatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter users by username", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		description := "Usuário administrador"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 2
			*args[1].(*string) = "admin_user"
			*args[2].(*string) = "admin@example.com"
			*args[3].(*string) = "hashed_password_admin"
			*args[4].(*string) = description
			*args[5].(*bool) = true
			*args[6].(*int) = 2
			*args[7].(*time.Time) = now.Add(-24 * time.Hour)
			*args[8].(*time.Time) = now
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "username",
				SortOrder: "desc",
			},
			Username: "admin",
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(2), result[0].UID)
		assert.Equal(t, "admin_user", result[0].Username)
		assert.Equal(t, "admin@example.com", result[0].Email)
		assert.Equal(t, "hashed_password_admin", result[0].Password)
		assert.Equal(t, description, result[0].Description)
		assert.True(t, result[0].Status)
		assert.Equal(t, 2, result[0].Version)
		assert.WithinDuration(t, now.Add(-24*time.Hour), result[0].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[0].UpdatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter users by email", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		description := "Usuário com email específico"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 3
			*args[1].(*string) = "jane_doe"
			*args[2].(*string) = "jane@example.com"
			*args[3].(*string) = "hashed_password_jane"
			*args[4].(*string) = description
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-48 * time.Hour)
			*args[8].(*time.Time) = now.Add(-24 * time.Hour)
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Email: "jane",
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(3), result[0].UID)
		assert.Equal(t, "jane_doe", result[0].Username)
		assert.Equal(t, "jane@example.com", result[0].Email)
		assert.Equal(t, "hashed_password_jane", result[0].Password)
		assert.Equal(t, description, result[0].Description)
		assert.True(t, result[0].Status)
		assert.Equal(t, 1, result[0].Version)
	})

	t.Run("successfully filter users by status and date range", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		createdFrom := now.Add(-7 * 24 * time.Hour)
		createdTo := now.Add(-1 * 24 * time.Hour)
		mockRows := new(mockDb.MockRows)

		// Primeiro usuário
		desc1 := "Usuário inativo 1"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 4
			*args[1].(*string) = "inactive_user1"
			*args[2].(*string) = "inactive1@example.com"
			*args[3].(*string) = "hashed_password_1"
			*args[4].(*string) = desc1
			*args[5].(*bool) = false
			*args[6].(*int) = 3
			*args[7].(*time.Time) = now.Add(-3 * 24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-2 * 24 * time.Hour)
		}).Return(nil).Once()

		// Segundo usuário
		desc2 := "Usuário inativo 2"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 5
			*args[1].(*string) = "inactive_user2"
			*args[2].(*string) = "inactive2@example.com"
			*args[3].(*string) = "hashed_password_2"
			*args[4].(*string) = desc2
			*args[5].(*bool) = false
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-5 * 24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-4 * 24 * time.Hour)
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		status := false
		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  20,
				Offset: 0,
			},
			Status:      &status,
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar primeiro usuário
		assert.Equal(t, int64(4), result[0].UID)
		assert.Equal(t, "inactive_user1", result[0].Username)
		assert.Equal(t, "inactive1@example.com", result[0].Email)
		assert.Equal(t, "hashed_password_1", result[0].Password)
		assert.Equal(t, desc1, result[0].Description)
		assert.False(t, result[0].Status)
		assert.Equal(t, 3, result[0].Version)

		// Verificar segundo usuário
		assert.Equal(t, int64(5), result[1].UID)
		assert.Equal(t, "inactive_user2", result[1].Username)
		assert.Equal(t, "inactive2@example.com", result[1].Email)
		assert.Equal(t, "hashed_password_2", result[1].Password)
		assert.Equal(t, desc2, result[1].Description)
		assert.False(t, result[1].Status)
		assert.Equal(t, 1, result[1].Version)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter users by updated_at range", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		updatedFrom := now.Add(-72 * time.Hour)
		updatedTo := now.Add(-24 * time.Hour)

		mockRows := new(mockDb.MockRows)

		// Primeiro usuário
		desc1 := "Atualizado há 2 dias"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 6
			*args[1].(*string) = "recent_user1"
			*args[2].(*string) = "recent1@example.com"
			*args[3].(*string) = "hashed_password_recent1"
			*args[4].(*string) = desc1
			*args[5].(*bool) = true
			*args[6].(*int) = 3
			*args[7].(*time.Time) = now.Add(-10 * 24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-48 * time.Hour)
		}).Return(nil).Once()

		// Segundo usuário
		desc2 := "Atualizado há 3 dias"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 7
			*args[1].(*string) = "recent_user2"
			*args[2].(*string) = "recent2@example.com"
			*args[3].(*string) = "hashed_password_recent2"
			*args[4].(*string) = desc2
			*args[5].(*bool) = false
			*args[6].(*int) = 2
			*args[7].(*time.Time) = now.Add(-15 * 24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-72 * time.Hour)
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "updated_at",
				SortOrder: "desc",
			},
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar ordenação por updated_at desc
		assert.Equal(t, int64(6), result[0].UID)
		assert.Equal(t, "recent_user1", result[0].Username)
		assert.WithinDuration(t, now.Add(-48*time.Hour), result[0].UpdatedAt, time.Second)

		assert.Equal(t, int64(7), result[1].UID)
		assert.Equal(t, "recent_user2", result[1].Username)
		assert.WithinDuration(t, now.Add(-72*time.Hour), result[1].UpdatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully filter users with invalid sort order defaults to asc", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		// Primeiro usuário (mais antigo)
		desc1 := "Mais antigo"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 8
			*args[1].(*string) = "user_a"
			*args[2].(*string) = "a@example.com"
			*args[3].(*string) = "hashed_password_a"
			*args[4].(*string) = desc1
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-48 * time.Hour)
			*args[8].(*time.Time) = now.Add(-48 * time.Hour)
		}).Return(nil).Once()

		// Segundo usuário (mais recente)
		desc2 := "Mais recente"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 9
			*args[1].(*string) = "user_b"
			*args[2].(*string) = "b@example.com"
			*args[3].(*string) = "hashed_password_b"
			*args[4].(*string) = desc2
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now.Add(-24 * time.Hour)
			*args[8].(*time.Time) = now.Add(-24 * time.Hour)
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "created_at",
				SortOrder: "INVALID_ORDER",
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar que ordenação é ASC (mais antigo primeiro)
		assert.Equal(t, int64(8), result[0].UID)
		assert.Equal(t, "user_a", result[0].Username)
		assert.WithinDuration(t, now.Add(-48*time.Hour), result[0].CreatedAt, time.Second)

		assert.Equal(t, int64(9), result[1].UID)
		assert.Equal(t, "user_b", result[1].Username)
		assert.WithinDuration(t, now.Add(-24*time.Hour), result[1].CreatedAt, time.Second)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully handle empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		emptyDescription := ""

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 10
			*args[1].(*string) = "no_description_user"
			*args[2].(*string) = "nodesc@example.com"
			*args[3].(*string) = "hashed_password_nodesc"
			*args[4].(*string) = emptyDescription
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = now
			*args[8].(*time.Time) = now
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(10), result[0].UID)
		assert.Equal(t, "no_description_user", result[0].Username)
		assert.Equal(t, "nodesc@example.com", result[0].Email)
		assert.Empty(t, result[0].Description)
	})

	t.Run("returns empty slice when no users found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Username: "nonexistent",
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result)

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("returns error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		expectedErr := fmt.Errorf("database connection failed")

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(nil, expectedErr)

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		result, err := repo.Filter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database connection failed")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("returns error when scanning row fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(fmt.Errorf("scan error: invalid data type"))

		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "scan error: invalid data type")
		assert.Contains(t, err.Error(), errMsg.ErrScan.Error())

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("returns error when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)

		// Primeira linha OK
		desc := "Usuário teste"

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 11
			*args[1].(*string) = "test_user"
			*args[2].(*string) = "test@example.com"
			*args[3].(*string) = "hashed_password_test"
			*args[4].(*string) = desc
			*args[5].(*bool) = true
			*args[6].(*int) = 1
			*args[7].(*time.Time) = time.Now()
			*args[8].(*time.Time) = time.Now()
		}).Return(nil).Once()

		// Erro na próxima chamada de Next
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(fmt.Errorf("cursor error: connection lost"))
		mockRows.On("Close").Return()

		filter := &filterUser.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.
			On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cursor error: connection lost")
		assert.Contains(t, err.Error(), errMsg.ErrIterate.Error())

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
