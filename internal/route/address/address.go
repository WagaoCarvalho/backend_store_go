package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/address/address"
	handlerFilter "github.com/WagaoCarvalho/backend_store_go/internal/handler/address/filter"
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
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	newRepoAddress := repoAddress.NewAddress(db)
	newRepoClientCpf := repoClientCpf.NewClientCpfRepo(db)
	newRepoUser := repoUser.NewUser(db)
	newRepoSupplier := repoSupplier.NewSupplier(db)

	newServiceAddress := service.NewAddressService(newRepoAddress, newRepoClientCpf, newRepoUser, newRepoSupplier)
	newHandlerAddress := handler.NewAddressHandler(newServiceAddress, log)

	newRepoFilter := repoFilter.NewFilterAddress(db)
	newServiceFilter := serviceFilter.NewAddressFilterService(newRepoFilter)
	newHandlerFilter := handlerFilter.NewAddressFilterHandler(newServiceFilter, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()

	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	const (
		addresses = "/addresses"
		enable    = "/enable"
		disable   = "/disable"
		filter    = "/filter"
	)

	s.HandleFunc(baseURL+addresses, newHandlerAddress.Create).Methods(http.MethodPost)

	s.HandleFunc(baseURL+addresses+idPath, newHandlerAddress.GetByID).Methods(http.MethodGet)

	s.HandleFunc(baseURL+addresses+idPath, newHandlerAddress.Update).Methods(http.MethodPut)
	s.HandleFunc(baseURL+addresses+idPath, newHandlerAddress.Delete).Methods(http.MethodDelete)

	s.HandleFunc(baseURL+addresses+idPath+enable, newHandlerAddress.Enable).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+addresses+idPath+disable, newHandlerAddress.Disable).Methods(http.MethodPatch)

	s.HandleFunc(baseURL+addresses+filter, newHandlerFilter.Filter).Methods(http.MethodGet)
}
