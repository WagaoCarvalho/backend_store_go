package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/address"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoClient "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	repoAddress := repoAddress.NewAddressRepository(db)
	repoClient := repoClient.NewClientRepository(db)
	repoUser := repoUser.NewUserRepository(db)
	repoSupplier := repoSupplier.NewSupplierRepository(db)
	service := service.NewAddressService(repoAddress, repoClient, repoUser, repoSupplier)
	handler := handler.NewAddressHandler(service, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/addresses", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/address/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/address/user/{id:[0-9]+}", handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc("/address/client/{id:[0-9]+}", handler.GetByClientID).Methods(http.MethodGet)
	s.HandleFunc("/address/supplier/{id:[0-9]+}", handler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/address/enable/{id:[0-9]+}", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc("/address/disable/{id:[0-9]+}", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/address/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
}
