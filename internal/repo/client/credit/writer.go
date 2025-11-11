package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/credit"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *clientCreditRepo) Create(ctx context.Context, credit *models.ClientCredit) (*models.ClientCredit, error) {
	const query = `
		INSERT INTO client_credits (client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at)
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

func (r *clientCreditRepo) Update(ctx context.Context, credit *models.ClientCredit) error {
	const query = `
		UPDATE client_credits
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

func (r *clientCreditRepo) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM client_credits WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
