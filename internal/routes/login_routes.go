package routes

import (
	"net/http"
	"time"

	jwt "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	login "github.com/WagaoCarvalho/backend_store_go/internal/auth/login"
	logout "github.com/WagaoCarvalho/backend_store_go/internal/auth/logout"
	pass "github.com/WagaoCarvalho/backend_store_go/internal/auth/password"
	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	logoutHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/logout"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter, blacklist logout.TokenBlacklist) {
	userRepo := userRepositories.NewUserRepository(db, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwt.NewJWTManager(jwtCfg.SecretKey, time.Hour)

	hasher := pass.BcryptHasher{}

	loginService := login.NewLoginService(userRepo, log, jwtManager, hasher)
	loginHandler := loginHandlers.NewLoginHandler(loginService, log)

	logoutService := logout.NewLogoutService(blacklist, log, jwtCfg.SecretKey)
	logoutHandler := logoutHandlers.NewLogoutHandler(logoutService, log)

	r.HandleFunc("/login", loginHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler.Logout).Methods(http.MethodPost)
}
