package routes

import (
	"net/http"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middlewares/jwt"
	repo_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/users"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo_user := repo_user.NewUserRepository(db, log)
	hasher := auth.BcryptHasher{}

	userService := services.NewUserService(repo_user, log, hasher)
	handler := handlers.NewUserHandler(userService, log)

	// Rota pública
	r.HandleFunc("/user", handler.Create).Methods(http.MethodPost)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, config.LoadConfig().Jwt.SecretKey)) // <- uso correto do middleware

	s.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user/id/{id}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/user/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/user/email/{email}", handler.GetByEmail).Methods(http.MethodGet)
	s.HandleFunc("/user/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user/enable/{id:[0-9]+}", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc("/user/disable/{id:[0-9]+}", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/user/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
