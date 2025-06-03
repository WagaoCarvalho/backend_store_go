package routes

import (
	"net/http"
	"time"

	loginServices "github.com/WagaoCarvalho/backend_store_go/internal/auth"
	loginHandlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	userRepositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"

	"github.com/WagaoCarvalho/backend_store_go/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterLoginRoutes(r *mux.Router, db *pgxpool.Pool) {
	userRepo := userRepositories.NewUserRepository(db)

	// 🔐 Carregar configuração JWT
	jwtCfg := config.LoadJwtConfig()

	// 🪙 Criar JWTManager
	jwtManager := loginServices.NewJWTManager(jwtCfg.SecretKey, time.Hour)

	// 💡 Injetar JWTManager no serviço
	service := loginServices.NewLoginService(userRepo, jwtManager)
	handler := loginHandlers.NewLoginHandler(service)

	r.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
}
