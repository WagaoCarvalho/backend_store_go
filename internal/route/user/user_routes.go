package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/user"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user"
	jwt_auth "github.com/WagaoCarvalho/backend_store_go/pkg/auth/jwt"
	auth "github.com/WagaoCarvalho/backend_store_go/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/pkg/middleware/jwt"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo_user := repo.NewUserRepository(db, log)
	hasher := auth.BcryptHasher{}

	userService := service.NewUserService(repo_user, log, hasher)
	handler := handler.NewUserHandler(userService, log)

	// Rota pública
	r.HandleFunc("/user", handler.Create).Methods(http.MethodPost)

	// Config JWT
	jwtCfg := config.LoadJwtConfig()
	jwtManager := jwt_auth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/users", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/user/id/{id}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/user/version/{id:[0-9]+}", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/user/email/{email}", handler.GetByEmail).Methods(http.MethodGet)
	s.HandleFunc("/user/name/{username}", handler.GetByName).Methods(http.MethodGet)
	s.HandleFunc("/user/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/user/enable/{id:[0-9]+}", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc("/user/disable/{id:[0-9]+}", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/user/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
