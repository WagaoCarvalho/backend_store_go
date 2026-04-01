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
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	newRepo := repo.NewContact(db)
	newService := service.NewContactService(newRepo)
	newHandler := handler.NewContactHandler(newService, log)

	jwtCfg := config.LoadJwtConfig()

	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddleware.IsAuthByBearerToken(blacklist, log, jwtManager))

	const (
		contact = "/contact"
	)

	s.HandleFunc(baseURL+contact, newHandler.Create).Methods(http.MethodPost)

	s.HandleFunc(baseURL+contact+idPath, newHandler.GetByID).Methods(http.MethodGet)

	s.HandleFunc(baseURL+contact+idPath, newHandler.Update).Methods(http.MethodPut)

	s.HandleFunc(baseURL+contact+idPath, newHandler.Delete).Methods(http.MethodDelete)
}
