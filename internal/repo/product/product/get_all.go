package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (r *productRepo) GetAll(ctx context.Context) ([]*models.Product, error) {
	const query = `
	SELECT id, supplier_id, product_name, manufacturer, product_description,
		cost_price, sale_price, stock_quantity, barcode,
		status, allow_discount, max_discount_percent,
		created_at, updated_at
	FROM products
	ORDER BY id;
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := scanProductRowLimited(rows, &p); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	return products, nil
}
