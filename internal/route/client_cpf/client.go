package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/client_cpf/client_cpf"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/handler/client_cpf/filter"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client_cpf/client"
	repoFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/client_cpf/filter"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client_cpf/client"
	serviceFilter "github.com/WagaoCarvalho/backend_store_go/internal/service/client_cpf/filter"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterClientRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	newRepo := repo.NewClientCpfRepo(db)
	newService := service.NewClientCpfService(newRepo)
	newHandler := handler.NewClientCpfHandler(newService, log)

	newRepoFilter := repoFilter.NewFilterClientCpf(db)
	newServiceFilter := serviceFilter.NewClientCpfFilterService(newRepoFilter)
	newFilter := filter.NewClientCpfFilterHandler(newServiceFilter, log)

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
		clients = "/clients-cpf"
		version = "/version"
		enable  = "/enable"
		disable = "/disable"
		filter  = "/filter"
	)

	s.HandleFunc(baseURL+clients, newHandler.Create).Methods(http.MethodPost)

	s.HandleFunc(baseURL+idPath+clients, newHandler.GetByID).Methods(http.MethodGet)

	s.HandleFunc(baseURL+idPath+clients+version, newHandler.GetVersionByID).Methods(http.MethodGet)

	s.HandleFunc(baseURL+idPath+clients, newHandler.Update).Methods(http.MethodPut)

	s.HandleFunc(baseURL+idPath+clients, newHandler.Delete).Methods(http.MethodDelete)

	s.HandleFunc(baseURL+idPath+clients+disable, newHandler.Disable).Methods(http.MethodPatch)

	s.HandleFunc(baseURL+idPath+clients+enable, newHandler.Enable).Methods(http.MethodPatch)

	s.HandleFunc(baseURL+clients+filter, newFilter.Filter).Methods(http.MethodGet)
}
