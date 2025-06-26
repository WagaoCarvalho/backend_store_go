package routes

import (
	"net/http"

	addressHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	addressRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	addressServices "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter) {
	repo := addressRepositories.NewAddressRepository(db, log)
	service := addressServices.NewAddressService(repo)
	handler := addressHandlers.NewAddressHandler(service)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken)

	s.HandleFunc("/addresses", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/address/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/address/user/{id:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/address/client/{id:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/address/supplier/{id:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
