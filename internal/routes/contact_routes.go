package routes

import (
	"net/http"

	contact_handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/contacts"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	contact_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	contact_services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterContactRoutes(r *mux.Router, db *pgxpool.Pool) {
	repo := contact_repositories.NewContactRepository(db)
	service := contact_services.NewContactService(repo)
	handler := contact_handlers.NewContactHandler(service)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken)

	s.HandleFunc("/contact", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/contact/user/{userID:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/contact/client/{clientID:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/contact/supplier/{supplierID:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
