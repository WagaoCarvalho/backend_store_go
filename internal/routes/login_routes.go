package routes

import (
	"net/http"

	jwt_auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	login "github.com/WagaoCarvalho/backend_store_go/internal/auth/login"
	logout "github.com/WagaoCarvalho/backend_store_go/internal/auth/logout"
	pass "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	login_handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	logout_handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/logout"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	user_repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/users"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist logout.TokenBlacklist,
) {
	userRepo := user_repositories.NewUserRepository(db, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Instanciar JWTManager que implementa JWTService
	jwtManager := jwt_auth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	hasher := pass.BcryptHasher{}

	loginService := login.NewLoginService(userRepo, log, jwtManager, hasher)
	loginHandler := login_handlers.NewLoginHandler(loginService, log)

	// Passa jwtManager (JWTService) em vez de string SecretKey
	logoutService := logout.NewLogoutService(blacklist, log, jwtManager)
	logoutHandler := logout_handlers.NewLogoutHandler(logoutService, log)

	r.HandleFunc("/login", loginHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler.Logout).Methods(http.MethodPost)
}
