package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/client_cpf"
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
	repo := repo.NewClientCpfRepo(db)
	service := service.NewClientCpfService(repo)
	handler := handler.NewClientCpfHandler(service, log)

	repoFilter := repoFilter.NewFilterClientCpf(db)
	serviceFilter := serviceFilter.NewClientCpfFilterService(repoFilter)
	filter := filter.NewClientCpfFilterHandler(serviceFilter, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/client-cpf", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/client-cpf/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/client-cpf/{id:[0-9]+}/version", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/client-cpf/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/client-cpf/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/client-cpf/{id:[0-9]+}/disable", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/client-cpf/{id:[0-9]+}/enable", handler.Enable).Methods(http.MethodPatch)

	s.HandleFunc("/clients-cpf/filter", filter.Filter).Methods(http.MethodGet)
}
