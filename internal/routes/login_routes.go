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

	// 🔐 Carregar configuração JWT
	jwtCfg := config.LoadJwtConfig()

	// 🪙 Criar JWTManager
	jwtManager := jwt.NewJWTManager(jwtCfg.SecretKey, time.Hour)

	// 🔑 Criar Hasher (bcrypt)
	hasher := pass.BcryptHasher{}

	// 💡 Injetar dependências no serviço de login
	loginService := login.NewLoginService(userRepo, log, jwtManager, hasher)
	loginHandler := loginHandlers.NewLoginHandler(loginService)

	// Serviço e handler de logout, injetando blacklist e secretKey
	logoutService := logout.NewLogoutService(blacklist, log, jwtCfg.SecretKey)
	logoutHandler := logoutHandlers.NewLogoutHandler(logoutService, log)

	r.HandleFunc("/login", loginHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler.Logout).Methods(http.MethodPost)
}
