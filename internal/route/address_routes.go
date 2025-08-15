package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/address"
	jwt_middlewares "github.com/WagaoCarvalho/backend_store_go/internal/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address"
	jwt_auth "github.com/WagaoCarvalho/backend_store_go/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt_middlewares.TokenBlacklist,
) {
	repo := repo.NewAddressRepository(db, log)
	service := service.NewAddressService(repo, log)
	handler := handler.NewAddressHandler(service, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Instanciar JWTManager que implementa JWTService
	jwtManager := jwt_auth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt_middlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/addresses", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/address/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/address/user/{id:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/address/client/{id:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/address/supplier/{id:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
