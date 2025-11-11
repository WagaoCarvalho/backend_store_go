package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/credit"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *clientCreditRepo) GetByID(ctx context.Context, id int64) (*models.ClientCredit, error) {
	const query = `
		SELECT id, client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at
		FROM client_credits
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

func (r *clientCreditRepo) GetByClientID(ctx context.Context, clientID int64) (*models.ClientCredit, error) {
	const query = `
		SELECT id, client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at
		FROM client_credits
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

func (r *clientCreditRepo) GetAll(ctx context.Context) ([]*models.ClientCredit, error) {
	const query = `
		SELECT id, client_id, allow_credit, credit_limit, credit_balance, created_at, updated_at
		FROM client_credits
		ORDER BY id;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var credits []*models.ClientCredit
	for rows.Next() {
		var credit models.ClientCredit
		if err := rows.Scan(
			&credit.ID,
			&credit.ClientID,
			&credit.AllowCredit,
			&credit.CreditLimit,
			&credit.CreditBalance,
			&credit.CreatedAt,
			&credit.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		credits = append(credits, &credit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return credits, nil
}

func (r *clientCreditRepo) GetByName(ctx context.Context, name string) ([]*models.ClientCredit, error) {
	const query = `
		SELECT cc.id, cc.client_id, cc.allow_credit, cc.credit_limit, cc.credit_balance, cc.created_at, cc.updated_at
		FROM client_credits cc
		INNER JOIN clients c ON cc.client_id = c.id
		WHERE c.name ILIKE '%' || $1 || '%'
		ORDER BY cc.id;
	`

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var credits []*models.ClientCredit
	for rows.Next() {
		var credit models.ClientCredit
		if err := rows.Scan(
			&credit.ID,
			&credit.ClientID,
			&credit.AllowCredit,
			&credit.CreditLimit,
			&credit.CreditBalance,
			&credit.CreatedAt,
			&credit.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		credits = append(credits, &credit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return credits, nil
}

func (r *clientCreditRepo) GetVersionByID(ctx context.Context, id int64) (int, error) {
	const query = `SELECT version FROM client_credits WHERE id = $1`
	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return version, nil
}
