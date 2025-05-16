package routes

import (
	"net/http"

	supplierCategoryHandler "github.com/WagaoCarvalho/backend_store_go/internal/handlers/supplier/supplier_categories"
	supplierCategoryRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_categories"
	supplierCategoryServices "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_categories"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierCategoryRoutes(r *mux.Router, db *pgxpool.Pool) {
	supplierCategoryRepo := supplierCategoryRepositories.NewSupplierCategoryRepository(db)
	supplierCategoryService := supplierCategoryServices.NewSupplierCategoryService(supplierCategoryRepo)
	supplierCategoryHandler := supplierCategoryHandler.NewSupplierCategoryHandler(supplierCategoryService)

	r.HandleFunc("/supplier-category", supplierCategoryHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/supplier-categories", supplierCategoryHandler.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.Delete).Methods(http.MethodDelete)
}
