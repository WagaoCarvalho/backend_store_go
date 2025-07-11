package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	middlewares "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist middlewares.TokenBlacklist, // <- injetar blacklist aqui
) {
	repo := repositories.NewAddressRepository(db, log)
	service := services.NewAddressService(repo, log)
	handler := handlers.NewAddressHandler(service, log)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken(blacklist, log, config.LoadConfig().Jwt.SecretKey)) // <- uso atualizado do middleware

	s.HandleFunc("/addresses", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/address/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/address/user/{id:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/address/client/{id:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/address/supplier/{id:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
