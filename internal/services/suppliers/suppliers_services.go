package services

import (
	"context"
	"errors"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers"
	services_address "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	services_contact "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	services_supplier_category "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_categories"
	services_supplier "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations"
)

// Erros padronizados
var (
	ErrSupplierNameRequired         = errors.New("nome do fornecedor é obrigatório")
	ErrSupplierCreateFailed         = errors.New("erro ao criar fornecedor")
	ErrRelationAlreadyExists        = errors.New("relação fornecedor-categoria já existe")
	ErrCheckRelationFailed          = errors.New("erro ao verificar existência da relação")
	ErrRelationCreateFailed         = errors.New("erro ao criar relação fornecedor-categoria")
	ErrCategoryFetchFailed          = errors.New("erro ao buscar categoria")
	ErrAddressCreateFailed          = errors.New("erro ao criar endereço")
	ErrContactCreateFailed          = errors.New("erro ao criar contato")
	ErrInvalidSupplierIDForDeletion = errors.New("ID inválido para deletar fornecedor")
	ErrInvalidSupplierID            = errors.New("ID do fornecedor é inválido")
	ErrSupplierVersionRequired      = errors.New("versão do fornecedor é obrigatória")

	// Regras de negócio
	ErrSupplierNotFound        = errors.New("fornecedor não encontrado")
	ErrSupplierVersionConflict = errors.New("conflito de versão ao atualizar o fornecedor")

	// Operações
	ErrSupplierUpdate = errors.New("erro ao atualizar fornecedor")
)

type SupplierService interface {
	Create(
		ctx context.Context,
		supplier *models.Supplier,
		categoryID int64,
		address *models_address.Address,
		contact *models_contact.Contact,
	) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
}

type supplierService struct {
	repo            repository.SupplierRepository
	relationService services_supplier.SupplierCategoryRelationService
	addressService  services_address.AddressService
	contactService  services_contact.ContactService
	categoryService services_supplier_category.SupplierCategoryService
}

func NewSupplierService(
	repo repository.SupplierRepository,
	relationService services_supplier.SupplierCategoryRelationService,
	addressService services_address.AddressService,
	contactService services_contact.ContactService,
	categoryService services_supplier_category.SupplierCategoryService,
) SupplierService {
	return &supplierService{
		repo:            repo,
		relationService: relationService,
		addressService:  addressService,
		contactService:  contactService,
		categoryService: categoryService,
	}
}

func (s *supplierService) Create(
	ctx context.Context,
	supplier *models.Supplier,
	categoryID int64,
	address *models_address.Address,
	contact *models_contact.Contact,
) (int64, error) {
	if supplier == nil || supplier.Name == "" {
		return 0, ErrSupplierNameRequired
	}

	createdSupplier, err := s.repo.Create(ctx, *supplier)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrSupplierCreateFailed, err)
	}
	supplier.ID = createdSupplier.ID

	if categoryID > 0 {
		exists, err := s.relationService.HasRelation(ctx, createdSupplier.ID, categoryID)
		if err != nil {
			return 0, fmt.Errorf("%w: %v", ErrCheckRelationFailed, err)
		}
		if exists {
			return 0, ErrRelationAlreadyExists
		}

		if _, err := s.relationService.Create(ctx, createdSupplier.ID, categoryID); err != nil {
			return 0, fmt.Errorf("%w: %v", ErrRelationCreateFailed, err)
		}

		category, err := s.categoryService.GetByID(ctx, categoryID)
		if err != nil {
			return 0, fmt.Errorf("%w: %v", ErrCategoryFetchFailed, err)
		}
		if category != nil {
			supplier.Categories = append(supplier.Categories, *category)
		}
	}

	if address != nil {
		address.SupplierID = &createdSupplier.ID
		createdAddress, err := s.addressService.Create(ctx, address)
		if err != nil {
			return 0, fmt.Errorf("%w: %v", ErrAddressCreateFailed, err)
		}
		supplier.Address = createdAddress
	}

	if contact != nil {
		contact.SupplierID = &createdSupplier.ID
		createdContact, err := s.contactService.Create(ctx, contact)
		if err != nil {
			return 0, fmt.Errorf("%w (fornecedor %d): %v", ErrContactCreateFailed, createdSupplier.ID, err)
		}
		supplier.Contact = createdContact
	}

	return createdSupplier.ID, nil
}

func (s *supplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *supplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	return s.repo.GetAll(ctx)
}

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) error {
	if supplier.ID <= 0 {
		return ErrInvalidSupplierID
	}

	if supplier.Name == "" {
		return ErrSupplierNameRequired
	}

	if supplier.Version == 0 {
		return ErrSupplierVersionRequired
	}

	err := s.repo.Update(ctx, supplier)
	if err != nil {
		if errors.Is(err, repository.ErrVersionConflict) {
			return ErrSupplierVersionConflict
		}
		if errors.Is(err, repository.ErrSupplierNotFound) {
			return ErrSupplierNotFound
		}
		return fmt.Errorf("%w: %v", ErrSupplierUpdate, err)
	}

	return nil
}

func (s *supplierService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrInvalidSupplierIDForDeletion
	}
	return s.repo.Delete(ctx, id)
}
