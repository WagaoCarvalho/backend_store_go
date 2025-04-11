package services

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"github.com/WagaoCarvalho/backend_store_go/internal/models"
// )

// type SupplierCategoryService interface {
// 	Create(ctx context.Context, category *models.SupplierCategory) (int64, error)
// 	GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error)
// 	GetAll(ctx context.Context) ([]*models.SupplierCategory, error)
// 	Update(ctx context.Context, category *models.SupplierCategory) error
// 	Delete(ctx context.Context, id int64) error
// }

// type supplierCategoryServiceImpl struct {
// 	repo repository.SupplierCategoryRepository
// }

// // Construtor
// func NewSupplierCategoryService(repo repository.SupplierCategoryRepository) SupplierCategoryService {
// 	return &supplierCategoryServiceImpl{repo: repo}
// }

// // Create valida e delega a criação ao repositório
// func (s *supplierCategoryServiceImpl) Create(ctx context.Context, category *models.SupplierCategory) (int64, error) {
// 	if strings.TrimSpace(category.Name) == "" {
// 		return 0, errors.New("nome da categoria é obrigatório")
// 	}
// 	return s.repo.Create(ctx, category)
// }

// // GetByID busca por ID
// func (s *supplierCategoryServiceImpl) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
// 	if id <= 0 {
// 		return nil, fmt.Errorf("ID inválido")
// 	}
// 	return s.repo.GetByID(ctx, id)
// }

// // GetAll retorna todas as categorias
// func (s *supplierCategoryServiceImpl) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
// 	return s.repo.GetAll(ctx)
// }

// // Update valida e atualiza categoria
// func (s *supplierCategoryServiceImpl) Update(ctx context.Context, category *models.SupplierCategory) error {
// 	if category.ID == 0 {
// 		return fmt.Errorf("ID da categoria é obrigatório")
// 	}
// 	if strings.TrimSpace(category.Name) == "" {
// 		return errors.New("nome da categoria é obrigatório")
// 	}
// 	return s.repo.Update(ctx, category)
// }

// // Delete valida e remove categoria
// func (s *supplierCategoryServiceImpl) Delete(ctx context.Context, id int64) error {
// 	if id <= 0 {
// 		return fmt.Errorf("ID inválido para exclusão")
// 	}
// 	return s.repo.Delete(ctx, id)
// }
