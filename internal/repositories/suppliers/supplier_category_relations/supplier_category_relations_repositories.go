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
	ErrRelationNotFound    = errors.New("relação supplier-categoria não encontrada")
	ErrRelationExists      = errors.New("relação já existe")
	ErrInvalidRelationData = errors.New("dados inválidos para relação")
)

type SupplierCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplier(ctx context.Context, supplierID int64) error
	CheckIfExists(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelationRepo struct {
	db *pgxpool.Pool
}

func NewSupplierCategoryRelationRepo(db *pgxpool.Pool) SupplierCategoryRelationRepository {
	return &supplierCategoryRelationRepo{db: db}
}

// Cria uma nova relação entre supplier e category
func (r *supplierCategoryRelationRepo) Create(ctx context.Context, relation *models.SupplierCategoryRelations) (*models.SupplierCategoryRelations, error) {
	if relation.SupplierID <= 0 || relation.CategoryID <= 0 {
		return nil, ErrInvalidRelationData
	}

	exists, err := r.CheckIfExists(ctx, relation.SupplierID, relation.CategoryID)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("erro ao criar relação: %w", err)
	}

	return relation, nil
}

// Verifica se já existe a relação
func (r *supplierCategoryRelationRepo) CheckIfExists(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	query := `
		SELECT 1 FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2 LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, supplierID, categoryID).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("erro ao verificar existência: %w", err)
	}

	return exists == 1, nil
}

// Retorna todas as relações para um supplier
func (r *supplierCategoryRelationRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelations, error) {
	query := `
		SELECT id, supplier_id, category_id, created_at, updated_at
		FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar relações do fornecedor: %w", err)
	}
	defer rows.Close()

	var relations []*models.SupplierCategoryRelations
	for rows.Next() {
		var rData models.SupplierCategoryRelations
		if err := rows.Scan(&rData.ID, &rData.SupplierID, &rData.CategoryID, &rData.CreatedAt, &rData.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao ler relação: %w", err)
		}
		relations = append(relations, &rData)
	}

	return relations, nil
}

// Retorna todas as relações para uma category
func (r *supplierCategoryRelationRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelations, error) {
	query := `
		SELECT id, supplier_id, category_id, created_at, updated_at
		FROM supplier_category_relations
		WHERE category_id = $1;
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar relações da categoria: %w", err)
	}
	defer rows.Close()

	var relations []*models.SupplierCategoryRelations
	for rows.Next() {
		var rData models.SupplierCategoryRelations
		if err := rows.Scan(&rData.ID, &rData.SupplierID, &rData.CategoryID, &rData.CreatedAt, &rData.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao ler relação: %w", err)
		}
		relations = append(relations, &rData)
	}

	return relations, nil
}

// Deleta uma única relação específica
func (r *supplierCategoryRelationRepo) Delete(ctx context.Context, supplierID, categoryID int64) error {
	query := `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1 AND category_id = $2;
	`

	cmd, err := r.db.Exec(ctx, query, supplierID, categoryID)
	if err != nil {
		return fmt.Errorf("erro ao deletar relação: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrRelationNotFound
	}

	return nil
}

// Deleta todas as relações de um fornecedor
func (r *supplierCategoryRelationRepo) DeleteAllBySupplier(ctx context.Context, supplierID int64) error {
	query := `
		DELETE FROM supplier_category_relations
		WHERE supplier_id = $1;
	`

	_, err := r.db.Exec(ctx, query, supplierID)
	if err != nil {
		return fmt.Errorf("erro ao deletar todas as relações do fornecedor: %w", err)
	}

	return nil
}
