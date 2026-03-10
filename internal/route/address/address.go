package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/address/address"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/handler/address/filter"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address/address"
	repoFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/address/filter"
	repoClientCpf "github.com/WagaoCarvalho/backend_store_go/internal/repo/client_cpf/client"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address/address"
	serviceFilter "github.com/WagaoCarvalho/backend_store_go/internal/service/address/filter"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAddressRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	// Load server configuration
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	repoAddress := repoAddress.NewAddress(db)
	repoClientCpf := repoClientCpf.NewClientCpfRepo(db)
	repoUser := repoUser.NewUser(db)
	repoSupplier := repoSupplier.NewSupplier(db)
	service := service.NewAddressService(repoAddress, repoClientCpf, repoUser, repoSupplier)
	handler := handler.NewAddressHandler(service, log)

	repoFilter := repoFilter.NewFilterAddress(db)
	serviceFilter := serviceFilter.NewAddressFilterService(repoFilter)
	filter := filter.NewAddressFilterHandler(serviceFilter, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()

	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	// Use the configuration values for routes
	s.HandleFunc(baseURL+"/addresses", handler.Create).Methods(http.MethodPost)

	s.HandleFunc(baseURL+idPath+"/addresses", handler.GetByID).Methods(http.MethodGet)

	s.HandleFunc(baseURL+idPath+"/addresses", handler.Update).Methods(http.MethodPut)
	s.HandleFunc(baseURL+idPath+"/addresses", handler.Delete).Methods(http.MethodDelete)

	s.HandleFunc(baseURL+idPath+"/addresses/enable", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+idPath+"/addresses/disable", handler.Disable).Methods(http.MethodPatch)

	s.HandleFunc(baseURL+"/addresses/filter", filter.Filter).Methods(http.MethodGet)
}
