package routes

import (
	"net/http"

	loginHandler "github.com/WagaoCarvalho/backend_store_go/internal/handler/login"
	logoutHandler "github.com/WagaoCarvalho/backend_store_go/internal/handler/logout"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	pass "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	login "github.com/WagaoCarvalho/backend_store_go/internal/service/login"
	logout "github.com/WagaoCarvalho/backend_store_go/internal/service/logout"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist logout.TokenBlacklist,
) {
	userRepo := repo.NewUserRepository(db, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Instanciar JWTManager que implementa JWTService
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	hasher := pass.BcryptHasher{}

	loginService := login.NewLoginService(userRepo, log, jwtManager, hasher)
	loginHandler := loginHandler.NewLoginHandler(loginService, log)

	// Passa jwtManager (JWTService) em vez de string SecretKey
	logoutService := logout.NewLogoutService(blacklist, log, jwtManager)
	logoutHandler := logoutHandler.NewLogoutHandler(logoutService, log)

	r.HandleFunc("/login", loginHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler.Logout).Methods(http.MethodPost)
}
