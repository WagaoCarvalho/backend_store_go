package repo

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	filterSale "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSale_Filter(t *testing.T) {
	t.Run("successfully get all sales", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		saleDate := now.Add(-24 * time.Hour)
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // id
			mock.AnythingOfType("**int64"),    // client_id
			mock.AnythingOfType("**int64"),    // user_id
			mock.AnythingOfType("*time.Time"), // sale_date
			mock.AnythingOfType("*float64"),   // total_items_amount
			mock.AnythingOfType("*float64"),   // total_items_discount
			mock.AnythingOfType("*float64"),   // total_sale_discount
			mock.AnythingOfType("*float64"),   // total_amount
			mock.AnythingOfType("*string"),    // payment_type
			mock.AnythingOfType("*string"),    // status
			mock.AnythingOfType("*string"),    // notes (string, não *string)
			mock.AnythingOfType("*int"),       // version
			mock.AnythingOfType("*time.Time"), // created_at
			mock.AnythingOfType("*time.Time"), // updated_at
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1

			clientID := int64(100)
			*args[1].(**int64) = &clientID

			userID := int64(50)
			*args[2].(**int64) = &userID

			*args[3].(*time.Time) = saleDate
			*args[4].(*float64) = 1000.0
			*args[5].(*float64) = 100.0
			*args[6].(*float64) = 50.0
			*args[7].(*float64) = 850.0
			*args[8].(*string) = "CREDIT_CARD"
			*args[9].(*string) = "COMPLETED"
			*args[10].(*string) = "Venda de teste" // string direta, não ponteiro
			*args[11].(*int) = 1
			*args[12].(*time.Time) = now
			*args[13].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.NotNil(t, result[0].ClientID)
		assert.Equal(t, int64(100), *result[0].ClientID)
		assert.NotNil(t, result[0].UserID)
		assert.Equal(t, int64(50), *result[0].UserID)
		assert.Equal(t, saleDate, result[0].SaleDate)
		assert.Equal(t, 1000.0, result[0].TotalItemsAmount)
		assert.Equal(t, 100.0, result[0].TotalItemsDiscount)
		assert.Equal(t, 50.0, result[0].TotalSaleDiscount)
		assert.Equal(t, 850.0, result[0].TotalAmount)
		assert.Equal(t, "CREDIT_CARD", result[0].PaymentType)
		assert.Equal(t, "COMPLETED", result[0].Status)
		assert.Equal(t, "Venda de teste", result[0].Notes) // string direta
		assert.Equal(t, 1, result[0].Version)
		assert.WithinDuration(t, now, result[0].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[0].UpdatedAt, time.Second)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("uses allowed sort field when SortBy is valid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			clientID := int64(100)
			*args[1].(**int64) = &clientID
			userID := int64(50)
			*args[2].(**int64) = &userID
			*args[3].(*time.Time) = now
			*args[4].(*float64) = 1000.0
			*args[11].(*int) = 1
			*args[12].(*time.Time) = now
			*args[13].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "total_amount",
				SortOrder: "desc",
			},
		}

		mockDB.
			On(
				"Query",
				ctx,
				mock.MatchedBy(func(q string) bool {
					return strings.Contains(q, "ORDER BY total_amount desc")
				}),
				mock.Anything,
			).
			Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("defaults sort order to asc when SortOrder is invalid", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		now := time.Now()

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			clientID := int64(100)
			*args[1].(**int64) = &clientID
			userID := int64(50)
			*args[2].(**int64) = &userID
			*args[3].(*time.Time) = now
			*args[11].(*int) = 1
			*args[12].(*time.Time) = now
			*args[13].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:     10,
				Offset:    0,
				SortBy:    "sale_date",
				SortOrder: "INVALID",
			},
		}

		mockDB.
			On(
				"Query",
				ctx,
				mock.MatchedBy(func(q string) bool {
					return strings.Contains(q, "ORDER BY sale_date asc")
				}),
				mock.Anything,
			).
			Return(mockRows, nil)

		_, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()
		dbErr := errors.New("database failure")

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		emptyRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(emptyRows, dbErr)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()
		scanErr := errors.New("failed to scan row")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()
		rowsErr := errors.New("iteration error")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters client_id and user_id correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 42
			clientID := int64(200)
			*args[1].(**int64) = &clientID
			userID := int64(75)
			*args[2].(**int64) = &userID
			*args[3].(*time.Time) = now
			*args[8].(*string) = "PIX"
			*args[9].(*string) = "PENDING"
			*args[10].(*string) = "Venda filtrada"
			*args[11].(*int) = 2
			*args[12].(*time.Time) = now
			*args[13].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		clientID := int64(200)
		userID := int64(75)
		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 5,
			},
			ClientID: &clientID,
			UserID:   &userID,
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			if len(args) != 2 {
				return false
			}
			cid := args[0].(int64)
			uid := args[1].(int64)
			return cid == 200 && uid == 75
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(42), result[0].ID)
		assert.Equal(t, int64(200), *result[0].ClientID)
		assert.Equal(t, int64(75), *result[0].UserID)
		assert.Equal(t, "PIX", result[0].PaymentType)
		assert.Equal(t, "PENDING", result[0].Status)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters status and payment_type correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 7
			clientID := int64(150)
			*args[1].(**int64) = &clientID
			userID := int64(60)
			*args[2].(**int64) = &userID
			*args[3].(*time.Time) = now
			*args[8].(*string) = "CASH"
			*args[9].(*string) = "COMPLETED"
			*args[10].(*string) = "Venda à vista"
			*args[11].(*int) = 3
			*args[12].(*time.Time) = now
			*args[13].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Status:      "COMPLETED",
			PaymentType: "CASH",
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			if len(args) != 2 {
				return false
			}
			status := args[0].(string)
			paymentType := args[1].(string)
			return status == "COMPLETED" && paymentType == "CASH"
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(7), result[0].ID)
		assert.Equal(t, "CASH", result[0].PaymentType)
		assert.Equal(t, "COMPLETED", result[0].Status)
		assert.Equal(t, 3, result[0].Version)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters date ranges correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &saleFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		saleDateFrom := now.Add(-72 * time.Hour)
		saleDateTo := now.Add(-24 * time.Hour)
		createdFrom := now.Add(-48 * time.Hour)
		createdTo := now.Add(-12 * time.Hour)

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("**int64"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*float64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 99
			clientID := int64(300)
			*args[1].(**int64) = &clientID
			userID := int64(80)
			*args[2].(**int64) = &userID
			*args[3].(*time.Time) = now.Add(-48 * time.Hour)
			*args[10].(*string) = "Filtro de data"
			*args[11].(*int) = 4
			*args[12].(*time.Time) = now.Add(-36 * time.Hour)
			*args[13].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterSale.SaleFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
			SaleDateFrom: &saleDateFrom,
			SaleDateTo:   &saleDateTo,
			CreatedFrom:  &createdFrom,
			CreatedTo:    &createdTo,
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			if len(args) != 4 {
				return false
			}
			hasSaleDateFrom := false
			hasSaleDateTo := false
			hasCreatedFrom := false
			hasCreatedTo := false

			for _, arg := range args {
				if t, ok := arg.(time.Time); ok {
					if t.Equal(saleDateFrom) {
						hasSaleDateFrom = true
					} else if t.Equal(saleDateTo) {
						hasSaleDateTo = true
					} else if t.Equal(createdFrom) {
						hasCreatedFrom = true
					} else if t.Equal(createdTo) {
						hasCreatedTo = true
					}
				}
			}

			return hasSaleDateFrom && hasSaleDateTo && hasCreatedFrom && hasCreatedTo
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(99), result[0].ID)
		assert.Equal(t, "Filtro de data", result[0].Notes)
		assert.Equal(t, 4, result[0].Version)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
