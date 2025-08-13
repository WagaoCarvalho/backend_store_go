package routes

import (
	"net/http"

	jwt_auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/config"
	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handler/address"
	jwt_middlewares "github.com/WagaoCarvalho/backend_store_go/internal/middleware/jwt"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/WagaoCarvalho/backend_store_go/logger"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt_middlewares.TokenBlacklist,
) {
	repo := repositories.NewAddressRepository(db, log)
	service := services.NewAddressService(repo, log)
	handler := handlers.NewAddressHandler(service, log)

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
