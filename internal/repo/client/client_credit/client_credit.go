package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client_credit"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientCreditRepository interface {
	Create(ctx context.Context, credit *models.ClientCredit) (*models.ClientCredit, error)
	GetByID(ctx context.Context, id int64) (*models.ClientCredit, error)
	GetByClientID(ctx context.Context, clientID int64) (*models.ClientCredit, error)
	Update(ctx context.Context, credit *models.ClientCredit) error
	Delete(ctx context.Context, id int64) error
}

type clientCreditRepository struct {
	db *pgxpool.Pool
}

func NewClientCreditRepository(db *pgxpool.Pool) ClientCreditRepository {
	return &clientCreditRepository{db: db}
}

func (r *clientCreditRepository) Create(ctx context.Context, credit *models.ClientCredit) (*models.ClientCredit, error) {
	const query = `
		INSERT INTO client_credit (client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		credit.ClientID,
		credit.AllowCredit,
		credit.CreditLimit,
		credit.CreditBalance,
	).Scan(&credit.ID, &credit.CreatedAt, &credit.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return credit, nil
}

func (r *clientCreditRepository) GetByID(ctx context.Context, id int64) (*models.ClientCredit, error) {
	const query = `
		SELECT id, client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at
		FROM client_credit
		WHERE id = $1
	`

	var credit models.ClientCredit
	err := r.db.QueryRow(ctx, query, id).Scan(
		&credit.ID,
		&credit.ClientID,
		&credit.AllowCredit,
		&credit.CreditLimit,
		&credit.CreditBalance,
		&credit.CreatedAt,
		&credit.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &credit, nil
}

func (r *clientCreditRepository) GetByClientID(ctx context.Context, clientID int64) (*models.ClientCredit, error) {
	const query = `
		SELECT id, client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at
		FROM client_credit
		WHERE client_id = $1
	`

	var credit models.ClientCredit
	err := r.db.QueryRow(ctx, query, clientID).Scan(
		&credit.ID,
		&credit.ClientID,
		&credit.AllowCredit,
		&credit.CreditLimit,
		&credit.CreditBalance,
		&credit.CreatedAt,
		&credit.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &credit, nil
}

func (r *clientCreditRepository) Update(ctx context.Context, credit *models.ClientCredit) error {
	const query = `
		UPDATE client_credit
		SET allow_credit = $1, credit_limit = $2, credit_balance = $3, updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(ctx, query,
		credit.AllowCredit,
		credit.CreditLimit,
		credit.CreditBalance,
		credit.ID,
	)

	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *clientCreditRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM client_credit WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
