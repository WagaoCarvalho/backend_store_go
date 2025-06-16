package routes

import (
	"net/http"

	contactHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/contacts"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	contactRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	contactServices "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterContactRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := contactRepositories.NewContactRepository(db)
	service := contactServices.NewContactService(repo)
	handler := contactHandlers.NewContactHandler(service)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken)

	s.HandleFunc("/contact", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/contact/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/contact/user/{userID:[0-9]+}", handler.GetByUser).Methods(http.MethodGet)
	s.HandleFunc("/contact/client/{clientID:[0-9]+}", handler.GetByClient).Methods(http.MethodGet)
	s.HandleFunc("/contact/supplier/{supplierID:[0-9]+}", handler.GetBySupplier).Methods(http.MethodGet)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
