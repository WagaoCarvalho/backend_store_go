package routes

import (
	"net/http"

	jwt_auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/user_full"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/middleware/jwt"
	repo_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	repo_contact "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repo_user_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	repo_user "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_full_repositories"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/users/user_full_services"
	"github.com/WagaoCarvalho/backend_store_go/logger"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserFullRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	repo_user := repo_user.NewUserRepository(db, log)
	repo_address := repo_address.NewAddressRepository(db, log)
	repo_contact := repo_contact.NewContactRepository(db, log)
	repo_user_cat_rel := repo_user_cat_rel.NewUserCategoryRelationRepositories(db, log)
	hasher := auth.BcryptHasher{}

	userService := services.NewUserFullService(repo_user, repo_address, repo_contact, repo_user_cat_rel, log, hasher)
	handler := handlers.NewUserFullHandler(userService, log)

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
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager)) // <- passa jwtManager, não string SecretKey

	s.HandleFunc("/user-full", handler.CreateFull).Methods(http.MethodPost)
}
