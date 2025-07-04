package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/auth"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter) {
	userRepo := repositories.NewUserRepository(db, log)

	// Instancia o hasher bcrypt
	hasher := auth.BcryptHasher{}

	userService := services.NewUserService(userRepo, log, hasher)
	handler := handlers.NewUserHandler(userService)

	r.HandleFunc("/user", handler.Create).Methods(http.MethodPost)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken)

	s.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user/id/{id}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/user/version/{id}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/user/email/{email}", handler.GetByEmail).Methods(http.MethodGet)
	s.HandleFunc("/user/{id}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user/{id}", handler.Delete).Methods(http.MethodDelete)
}
