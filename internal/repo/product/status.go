package repo

import (
	"context"
	"fmt"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

// EnableProduct ativa um produto, alterando o campo `status` para TRUE.
func (r *productRepository) EnableProduct(ctx context.Context, uid int64) error {
	const query = `
		UPDATE products
		SET status = TRUE, updated_at = NOW()
		WHERE id = $1;
	`

	cmd, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	if cmd.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

// DisableProduct desativa um produto, alterando o campo `status` para FALSE.
func (r *productRepository) DisableProduct(ctx context.Context, uid int64) error {
	const query = `
		UPDATE products
		SET status = FALSE, updated_at = NOW()
		WHERE id = $1;
	`

	cmd, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	if cmd.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
