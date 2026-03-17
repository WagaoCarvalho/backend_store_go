package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	loginHandler "github.com/WagaoCarvalho/backend_store_go/internal/handler/login"
	logoutHandler "github.com/WagaoCarvalho/backend_store_go/internal/handler/logout"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	pass "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/filter"
	login "github.com/WagaoCarvalho/backend_store_go/internal/service/login"
	logout "github.com/WagaoCarvalho/backend_store_go/internal/service/logout"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist logout.TokenBlacklist,
) {
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL

	userRepo := repo.NewUserFilter(db)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	hasher := pass.BcryptHasher{}

	newLoginService := login.NewLoginService(userRepo, jwtManager, hasher)
	newLoginHandler := loginHandler.NewLoginHandler(newLoginService, log)

	newLogoutService := logout.NewLogoutService(blacklist, jwtManager)
	newLogoutHandler := logoutHandler.NewLogoutHandler(newLogoutService, log)

	s := r.PathPrefix("/").Subrouter()

	const (
		loginPath  = "/login"
		logoutPath = "/logout"
	)

	s.HandleFunc(baseURL+loginPath, newLoginHandler.Login).Methods(http.MethodPost)
	s.HandleFunc(baseURL+logoutPath, newLogoutHandler.Logout).Methods(http.MethodPost)
}
