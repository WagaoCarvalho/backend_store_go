package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/contacts"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterContactRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist, // <- injeta blacklist aqui
) {
	repo := repositories.NewContactRepository(db, log)
	service := services.NewContactService(repo, log)
	handler := handlers.NewContactHandler(service, log)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, config.LoadConfig().Jwt.SecretKey)) // <- uso do middleware com blacklist

	s.HandleFunc("/contact", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/contact/user/{userID:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/contact/client/{clientID:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/contact/supplier/{supplierID:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
