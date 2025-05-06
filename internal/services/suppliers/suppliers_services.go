package services

import (
	"context"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers"
	services_address "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	services_contact "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	services_supplier "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations"
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
}

func NewSupplierService(
	repo repository.SupplierRepository,
	relationService services_supplier.SupplierCategoryRelationService,
	addressService services_address.AddressService,
	contactService services_contact.ContactService,
) SupplierService {
	return &supplierService{
		repo:            repo,
		relationService: relationService,
		addressService:  addressService,
		contactService:  contactService,
	}
}

func (s *supplierService) Create(
	ctx context.Context,
	supplier *models.Supplier,
	categoryID int64,
	address *models_address.Address,
	contact *models_contact.Contact,
) (int64, error) {
	if supplier.Name == "" {
		return 0, fmt.Errorf("nome do fornecedor é obrigatório")
	}

	supplierID, err := s.repo.Create(ctx, supplier)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar fornecedor: %w", err)
	}
	supplier.ID = supplierID

	if categoryID > 0 {
		exists, err := s.relationService.HasRelation(ctx, supplierID, categoryID)
		if err != nil {
			return 0, fmt.Errorf("erro ao verificar existência da relação: %w", err)
		}
		if exists {
			return 0, fmt.Errorf("relação fornecedor-categoria já existe")
		}

		_, err = s.relationService.Create(ctx, supplierID, categoryID)
		if err != nil {
			return 0, fmt.Errorf("erro ao criar relação fornecedor-categoria: %w", err)
		}
	}

	if address != nil {
		address.SupplierID = supplierID

		_, err := s.addressService.Create(ctx, *address)
		if err != nil {
			return 0, fmt.Errorf("erro ao criar endereço: %w", err)
		}
	}

	if contact != nil {
		contact.SupplierID = &supplierID

		err := s.contactService.Create(ctx, contact)
		if err != nil {
			return 0, fmt.Errorf("erro ao criar contato: %w", err)
		}
	}

	return supplierID, nil
}

func (s *supplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *supplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	return s.repo.GetAll(ctx)
}

func (s *supplierService) Update(ctx context.Context, supplier *models.Supplier) error {
	if supplier.Name == "" {
		return fmt.Errorf("nome do fornecedor é obrigatório")
	}
	return s.repo.Update(ctx, supplier)
}

func (s *supplierService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return fmt.Errorf("ID inválido para deletar fornecedor")
	}
	return s.repo.Delete(ctx, id)
}
