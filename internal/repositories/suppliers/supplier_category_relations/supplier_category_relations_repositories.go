package repository

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRelationNotFound                        = errors.New("relação supplier-categoria não encontrada")
	ErrRelationExists                          = errors.New("relação já existe")
	ErrInvalidRelationData                     = errors.New("dados inválidos para relação")
	ErrCreateRelation                          = errors.New("erro ao criar relação")
	ErrCheckRelation                           = errors.New("erro ao verificar existência da relação")
	ErrGetRelationsBySupplier                  = errors.New("erro ao buscar relações do fornecedor")
	ErrGetRelationsByCategory                  = errors.New("erro ao buscar relações da categoria")
	ErrScanRelationRow                         = errors.New("erro ao ler relação")
	ErrDeleteRelation                          = errors.New("erro ao deletar relação")
	ErrDeleteAllRelationsBySupplier            = errors.New("erro ao deletar todas as relações do fornecedor")
	ErrInvalidSupplierCategoryRelationID       = errors.New("ID da relação de categoria do fornecedor é inválido")
	ErrSupplierCategoryRelationVersionRequired = errors.New("versão da relação de categoria do fornecedor é obrigatória")
	ErrSupplierCategoryRelationNotFound        = errors.New("relação de categoria do fornecedor não encontrada")
	ErrSupplierCategoryRelationUpdate          = errors.New("erro ao atualizar a relação de categoria do fornecedor")
)

type SupplierCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierId(ctx context.Context, supplierID int64) error
	HasSupplierCategoryRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
	Update(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
}

type supplierCategoryRelationRepo struct {
	db *pgxpool.Pool
}

func NewSupplierCategoryRelationRepo(db *pgxpool.Pool) SupplierCategoryRelationRepository {
	return &supplierCategoryRelationRepo{db: db}
}

func (r *supplierCategoryRelationRepo) Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	if relation.SupplierID <= 0 || relation.CategoryID <= 0 {
		return nil, ErrInvalidRelationData
	}

	exists, err := r.HasSupplierCategoryRelation(ctx, relation.SupplierID, relation.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCheckRelation, err)
	}
	if exists {
		return nil, ErrRelationExists
	}

	query := `
		INSERT INTO supplier_category_relations (supplier_id, category_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING created_at, updated_at;
	`

	err = r.db.QueryRow(ctx, query, relation.SupplierID, relation.CategoryID).
		Scan(&relation.CreatedAt, &relation.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	return relation, nil
}

func (r *supplierCategoryRelationRepo) HasSupplierCategoryRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	query := `
		SELECT 1 FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2 LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, supplierID, categoryID).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("%w: %v", ErrCheckRelation, err)
	}

	return exists == 1, nil
}

func (r *supplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	query := `
		SELECT id, supplier_id, category_id, created_at, updated_at
		FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsBySupplier, err)
	}
	defer rows.Close()

	var relations []*models.SupplierCategoryRelations
	for rows.Next() {
		var rData models.SupplierCategoryRelations
		if err := rows.Scan(&rData.ID, &rData.SupplierID, &rData.CategoryID, &rData.CreatedAt, &rData.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanRelationRow, err)
		}
		relations = append(relations, &rData)
	}

	return relations, nil
}

func (r *supplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	query := `
		SELECT id, supplier_id, category_id, created_at, updated_at
		FROM supplier_category_relations
		WHERE category_id = $1;
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsByCategory, err)
	}
	defer rows.Close()

	var relations []*models.SupplierCategoryRelations
	for rows.Next() {
		var rData models.SupplierCategoryRelations
		if err := rows.Scan(&rData.ID, &rData.SupplierID, &rData.CategoryID, &rData.CreatedAt, &rData.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanRelationRow, err)
		}
		relations = append(relations, &rData)
	}

	return relations, nil
}

func (r *supplierCategoryRelationRepo) Update(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	if relation.ID <= 0 {
		return nil, ErrInvalidSupplierCategoryRelationID
	}

	if relation.Version <= 0 {
		return nil, ErrSupplierCategoryRelationVersionRequired
	}

	query := `
		UPDATE supplier_category_relations
		SET supplier_id = $1,
		    category_id = $2,
		    updated_at = NOW(),
		    version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING id, supplier_id, category_id, created_at, updated_at, version
	`

	row := r.db.QueryRow(ctx, query,
		relation.SupplierID,
		relation.CategoryID,
		relation.ID,
		relation.Version,
	)

	updated := &models.SupplierCategoryRelations{}
	err := row.Scan(
		&updated.ID,
		&updated.SupplierID,
		&updated.CategoryID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.Version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSupplierCategoryRelationNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrSupplierCategoryRelationUpdate, err)
	}

	return updated, nil
}

func (r *supplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	query := `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2;
	`

	cmd, err := r.db.Exec(ctx, query, supplierID, categoryID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrRelationNotFound
	}

	return nil
}

func (r *supplierCategoryRelationRepo) DeleteAllBySupplierId(ctx context.Context, supplierID int64) error {
	query := `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	_, err := r.db.Exec(ctx, query, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAllRelationsBySupplier, err)
	}

	return nil
}
