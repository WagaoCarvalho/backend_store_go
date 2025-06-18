package routes

import (
	"net/http"

	addressHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	addressRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	addressServices "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := addressRepositories.NewAddressRepository(db)
	service := addressServices.NewAddressService(repo)
	handler := addressHandlers.NewAddressHandler(service)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken)

	s.HandleFunc("/addresses", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/address/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/address/user/{id:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/address/client/{id:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/address/supplier/{id:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/address/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
}
