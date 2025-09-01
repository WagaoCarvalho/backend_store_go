package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/contact"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddleware "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/contact"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterContactRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddleware.TokenBlacklist,
) {
	repo := repo.NewContactRepository(db, log)
	service := service.NewContactService(repo, log)
	handler := handler.NewContactHandler(service, log)

	// Carrega a configuração do JWT
	jwtCfg := config.LoadJwtConfig()

	// Instancia um JWTManager (implementa JWTService)
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddleware.IsAuthByBearerToken(blacklist, log, jwtManager)) // <- agora com JWTService válido

	s.HandleFunc("/contact", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/contact/user/{userID:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/contact/client/{clientID:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/contact/supplier/{supplierID:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/contact/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
