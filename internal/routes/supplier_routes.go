package routes

import (
	"net/http"

	supplierHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/supplier"
	supplierCategoryHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/supplier/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	addressRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	contactRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	supplierRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers"
	supplierCategoryRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_categories"
	supplierCategoryRelationRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
	addressServices "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	contactServices "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	supplierServices "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers"
	supplierCategoryServices "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_categories"
	supplierCategoryRelations "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter) {
	supplierRepo := supplierRepositories.NewSupplierRepository(db)
	supplierCategoryRepo := supplierCategoryRepositories.NewSupplierCategoryRepository(db, log)
	supplierCategoryRelationRepo := supplierCategoryRelationRepositories.NewSupplierCategoryRelationRepo(db)
	addressRepo := addressRepositories.NewAddressRepository(db, log)
	contactRepo := contactRepositories.NewContactRepository(db, log)

	relationService := supplierCategoryRelations.NewSupplierCategoryRelationService(supplierCategoryRelationRepo)
	categoryService := supplierCategoryServices.NewSupplierCategoryService(supplierCategoryRepo)
	addressService := addressServices.NewAddressService(addressRepo, log)
	contactService := contactServices.NewContactService(contactRepo, log)
	supplierService := supplierServices.NewSupplierService(supplierRepo, relationService, addressService, contactService, categoryService)

	supplierHandler := supplierHandlers.NewSupplierHandler(supplierService)
	categoryHandler := supplierCategoryHandlers.NewSupplierCategoryHandler(categoryService)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken)

	s.HandleFunc("/supplier", supplierHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/supplier/{id:[0-9]+}", supplierHandler.GetSupplierByID).Methods(http.MethodGet)
	s.HandleFunc("/suppliers", supplierHandler.GetAllSuppliers).Methods(http.MethodGet)
	s.HandleFunc("/supplier/{id:[0-9]+}", supplierHandler.UpdateSupplier).Methods(http.MethodPut)
	s.HandleFunc("/supplier/{id:[0-9]+}", supplierHandler.DeleteSupplier).Methods(http.MethodDelete)

	s.HandleFunc("/supplier-category", categoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", categoryHandler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/supplier-categories", categoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", categoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", categoryHandler.Delete).Methods(http.MethodDelete)
}
