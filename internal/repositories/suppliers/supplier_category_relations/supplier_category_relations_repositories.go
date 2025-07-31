package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierId(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelationRepo struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewSupplierCategoryRelationRepo(db *pgxpool.Pool, lg logger.LoggerAdapterInterface) SupplierCategoryRelationRepository {
	return &supplierCategoryRelationRepo{db: db, logger: lg}
}

func (r *supplierCategoryRelationRepo) Create(ctx context.Context, rel *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	const query = `
		INSERT INTO supplier_category_relations (supplier_id, category_id, created_at)
		VALUES ($1, $2, NOW())
		RETURNING created_at;
	`

	err := r.db.QueryRow(ctx, query, rel.SupplierID, rel.CategoryID).Scan(&rel.CreatedAt)
	if err != nil {
		r.logger.Error(ctx, err, "[Create] falha ao inserir relação", map[string]any{
			"supplier_id": rel.SupplierID,
			"category_id": rel.CategoryID,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	r.logger.Info(ctx, "[Create] relação criada com sucesso", map[string]any{
		"supplier_id": rel.SupplierID,
		"category_id": rel.CategoryID,
	})
	return rel, nil
}

func (r *supplierCategoryRelationRepo) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	const query = `
		SELECT 1 FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2 LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, supplierID, categoryID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		r.logger.Error(ctx, err, "[HasRelation] erro ao verificar existência", map[string]any{
			"supplier_id": supplierID, "category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", ErrCheckRelation, err)
	}
	return true, nil
}

func (r *supplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	const query = `
		SELECT supplier_id, category_id, created_at
		FROM supplier_category_relations
		WHERE supplier_id = $1
		ORDER BY category_id;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		r.logger.Error(ctx, err, "[GetBySupplierID] erro na consulta", map[string]any{"supplier_id": supplierID})
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsBySupplier, err)
	}
	defer rows.Close()

	var rels []*models.SupplierCategoryRelations
	for rows.Next() {
		rel := new(models.SupplierCategoryRelations)
		if err := rows.Scan(&rel.SupplierID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			r.logger.Error(ctx, err, "[GetBySupplierID] erro ao escanear linha", map[string]any{"supplier_id": supplierID})
			return nil, fmt.Errorf("%w: %v", ErrScanRelationRow, err)
		}
		rels = append(rels, rel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "[GetBySupplierID] erro de iteração", map[string]any{"supplier_id": supplierID})
		return nil, fmt.Errorf("%w: %v", ErrScanRelationRow, err)
	}
	return rels, nil
}

func (r *supplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	const query = `
		SELECT supplier_id, category_id, created_at
		FROM supplier_category_relations
		WHERE category_id = $1
		ORDER BY supplier_id;
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		r.logger.Error(ctx, err, "[GetByCategoryID] erro na consulta", map[string]any{"category_id": categoryID})
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsByCategory, err)
	}
	defer rows.Close()

	var rels []*models.SupplierCategoryRelations
	for rows.Next() {
		rel := new(models.SupplierCategoryRelations)
		if err := rows.Scan(&rel.SupplierID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			r.logger.Error(ctx, err, "[GetByCategoryID] erro ao escanear linha", map[string]any{"category_id": categoryID})
			return nil, fmt.Errorf("%w: %v", ErrScanRelationRow, err)
		}
		rels = append(rels, rel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "[GetByCategoryID] erro de iteração", map[string]any{"category_id": categoryID})
		return nil, fmt.Errorf("%w: %v", ErrScanRelationRow, err)
	}
	return rels, nil
}

func (r *supplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	const query = `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2;
	`

	cmd, err := r.db.Exec(ctx, query, supplierID, categoryID)
	if err != nil {
		r.logger.Error(ctx, err, "[Delete] erro ao deletar relação", map[string]any{
			"supplier_id": supplierID, "category_id": categoryID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrRelationNotFound
	}
	return nil
}

func (r *supplierCategoryRelationRepo) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	const query = `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	_, err := r.db.Exec(ctx, query, supplierID)
	if err != nil {
		r.logger.Error(ctx, err, "[DeleteAllBySupplierId] erro ao deletar todas relações", map[string]any{
			"supplier_id": supplierID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAllRelationsBySupplier, err)
	}
	return nil
}
