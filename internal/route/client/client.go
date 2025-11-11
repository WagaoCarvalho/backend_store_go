package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/client"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client/client"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterClientRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	repo := repo.NewClient(db)
	service := service.NewClientService(repo)
	handler := handler.NewClientHandler(service, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/client", handler.Create).Methods(http.MethodPost)
	s.HandleFunc("/client/{id:[0-9]+}", handler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/clients/name/{name}", handler.GetByName).Methods(http.MethodGet)
	s.HandleFunc("/client/{id:[0-9]+}/version", handler.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc("/clients/filter", handler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/client/{id:[0-9]+}", handler.Update).Methods(http.MethodPut)
	s.HandleFunc("/client/{id:[0-9]+}", handler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/client/{id:[0-9]+}/disable", handler.Disable).Methods(http.MethodPatch)
	s.HandleFunc("/client/{id:[0-9]+}/enable", handler.Enable).Methods(http.MethodPatch)
	s.HandleFunc("/client/{id:[0-9]+}/exists", handler.ClientExists).Methods(http.MethodGet)

}
