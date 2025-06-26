package routes

import (
	"net/http"

	user_handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	address_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	contact_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	user_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	user_category_relation_repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	user_services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter) {
	userRepo := user_repositories.NewUserRepository(db)
	addressRepo := address_repositories.NewAddressRepository(db, log)
	contactRepo := contact_repositories.NewContactRepository(db)
	relationRepo := user_category_relation_repo.NewUserCategoryRelationRepositories(db)

	userService := user_services.NewUserService(userRepo, addressRepo, contactRepo, relationRepo)
	handler := user_handlers.NewUserHandler(userService)

	r.HandleFunc("/user", handler.Create).Methods(http.MethodPost)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken)

	//s.HandleFunc("/user", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user/id/{id}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/user/version/{id}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/user/email/{email}", handler.GetByEmail).Methods(http.MethodGet)
	s.HandleFunc("/user/{id}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user/{id}", handler.Delete).Methods(http.MethodDelete)
}
