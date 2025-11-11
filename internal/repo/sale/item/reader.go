package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *saleItemRepo) GetByID(ctx context.Context, id int64) (*models.SaleItem, error) {
	const query = `
		SELECT 
			id, sale_id, product_id, quantity, unit_price, discount, tax,
			subtotal, description, created_at, updated_at
		FROM sale_items
		WHERE id = $1;
	`

	var item models.SaleItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.SaleID,
		&item.ProductID,
		&item.Quantity,
		&item.UnitPrice,
		&item.Discount,
		&item.Tax,
		&item.Subtotal,
		&item.Description,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &item, nil
}

func (r *saleItemRepo) GetBySaleID(ctx context.Context, saleID int64, limit, offset int) ([]*models.SaleItem, error) {
	const query = `
		SELECT 
			id, sale_id, product_id, quantity, unit_price, discount, tax,
			subtotal, description, created_at, updated_at
		FROM sale_items
		WHERE sale_id = $1
		ORDER BY id ASC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(ctx, query, saleID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var items []*models.SaleItem
	for rows.Next() {
		var item models.SaleItem
		if err := rows.Scan(
			&item.ID,
			&item.SaleID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Discount,
			&item.Tax,
			&item.Subtotal,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		items = append(items, &item)
	}

	return items, nil
}

func (r *saleItemRepo) GetByProductID(ctx context.Context, productID int64, limit, offset int) ([]*models.SaleItem, error) {
	const query = `
		SELECT 
			id, sale_id, product_id, quantity, unit_price, discount, tax,
			subtotal, description, created_at, updated_at
		FROM sale_items
		WHERE product_id = $1
		ORDER BY id ASC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.Query(ctx, query, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var items []*models.SaleItem
	for rows.Next() {
		var item models.SaleItem
		if err := rows.Scan(
			&item.ID,
			&item.SaleID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Discount,
			&item.Tax,
			&item.Subtotal,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		items = append(items, &item)
	}

	return items, nil
}
