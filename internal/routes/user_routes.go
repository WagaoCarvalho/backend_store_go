package routes

import (
	"net/http"

	userHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/middlewares"
	addressRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	contactRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	userCategoryRelationRepo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	userServices "github.com/WagaoCarvalho/backend_store_go/internal/services/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool) {
	userRepo := userRepositories.NewUserRepository(db)
	addressRepo := addressRepositories.NewAddressRepository(db)
	contactRepo := contactRepositories.NewContactRepository(db)
	relationRepo := userCategoryRelationRepo.NewUserCategoryRelationRepositories(db)

	userService := userServices.NewUserService(userRepo, addressRepo, contactRepo, relationRepo)
	handler := userHandlers.NewUserHandler(userService)

	r.HandleFunc("/user", handler.Create).Methods(http.MethodPost)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middlewares.IsAuthByBearerToken)

	//s.HandleFunc("/user", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user/id/{id}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/user/version/{id}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/user/email/{email}", handler.GetByEmail).Methods(http.MethodGet)
	s.HandleFunc("/user/{id}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user/{id}", handler.Delete).Methods(http.MethodDelete)
}
