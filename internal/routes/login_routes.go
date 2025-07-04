package routes

import (
	"net/http"
	"time"

	loginServices "github.com/WagaoCarvalho/backend_store_go/internal/auth"
	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"

	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(r *mux.Router, db *pgxpool.Pool, log *logger.LoggerAdapter) {
	userRepo := userRepositories.NewUserRepository(db, log)

	// 🔐 Carregar configuração JWT
	jwtCfg := config.LoadJwtConfig()

	// 🪙 Criar JWTManager
	jwtManager := loginServices.NewJWTManager(jwtCfg.SecretKey, time.Hour)

	// 🔑 Criar Hasher (bcrypt)
	hasher := loginServices.BcryptHasher{}

	// 💡 Injetar dependências no serviço de login
	service := loginServices.NewLoginService(userRepo, jwtManager, hasher)

	handler := loginHandlers.NewLoginHandler(service)

	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
}
