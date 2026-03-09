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

	const (
		baseUrl   = "/api/v1"
		idPath    = "/{id:[0-9]+}"
		addresses = "/addresses"
	)

	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc(baseUrl+addresses, handler.Create).Methods(http.MethodPost)

	s.HandleFunc(baseUrl+idPath+addresses, handler.GetByID).Methods(http.MethodGet)

	s.HandleFunc(baseUrl+idPath+addresses, handler.GetByUserID).Methods(http.MethodGet)
	s.HandleFunc(baseUrl+idPath+addresses, handler.GetByClientCpfID).Methods(http.MethodGet)
	s.HandleFunc(baseUrl+idPath+addresses, handler.GetBySupplierID).Methods(http.MethodGet)

	s.HandleFunc(baseUrl+idPath+addresses, handler.Update).Methods(http.MethodPut)
	s.HandleFunc(baseUrl+idPath+addresses, handler.Delete).Methods(http.MethodDelete)

	s.HandleFunc(baseUrl+idPath+addresses+"/enable", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc(baseUrl+idPath+addresses+"/disable", handler.Disable).Methods(http.MethodPatch)

	s.HandleFunc(baseUrl+addresses+"/filter", filter.Filter).Methods(http.MethodGet)

}
