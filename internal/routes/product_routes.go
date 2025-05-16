package routes

import (
	"net/http"

	productHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/product"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	productRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
	productServices "github.com/WagaoCarvalho/backend_store_go/internal/services/products"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := productRepositories.NewProductRepository(db)
	service := productServices.NewProductService(repo)
	handler := productHandlers.NewProductHandler(service)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken)

	s.HandleFunc("/products", handler.GetProducts).Methods(http.MethodGet)
	s.HandleFunc("/product/{id}", handler.GetProductById).Methods(http.MethodGet)
	s.HandleFunc("/product", handler.CreateProduct).Methods(http.MethodPost)
	s.HandleFunc("/product/{id}", handler.UpdateProduct).Methods(http.MethodPut)
	s.HandleFunc("/product/{id}", handler.DeleteProductById).Methods(http.MethodDelete)
	s.HandleFunc("/products/search", handler.GetProductsByName).Methods(http.MethodGet)
	s.HandleFunc("/products/price", handler.GetProductsBySalePriceRange).Methods(http.MethodGet)
	s.HandleFunc("/products/low-stock", handler.GetProductsLowInStock).Methods(http.MethodGet)
}
