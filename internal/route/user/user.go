package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handlerFilter "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/filter"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/user/user"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	auth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/password"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repoFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/filter"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	serviceFilter "github.com/WagaoCarvalho/backend_store_go/internal/service/user/filter"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	// Repositórios
	newRepoUser := repo.NewUser(db)
	newRepoFilter := repoFilter.NewUserFilter(db)

	// Dependências
	hasher := auth.BcryptHasher{}

	// Serviços
	newServiceUser := service.NewUserService(newRepoUser, hasher)
	newServiceFilter := serviceFilter.NewUserFilterService(newRepoFilter)

	// Handlers
	newHandlerUser := handler.NewUserHandler(newServiceUser, log)
	newHandlerFilter := handlerFilter.NewUserFilterHandler(newServiceFilter, log)

	// Config JWT
	jwtCfg := config.LoadJwtConfig()
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	// Rota pública (fora do subrouter protegido)
	r.HandleFunc(baseURL+"/user", newHandlerUser.Create).Methods(http.MethodPost)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	// Constantes para caminhos
	const (
		users   = "/users"
		user    = "/user"
		id      = "/id/"
		version = "/version/"
		email   = "/email/"
		name    = "/name/"
		enable  = "/enable/"
		disable = "/disable/"
		filter  = "/filter"
		idParam = "{id:[0-9]+}"
	)

	// Rotas de listagem e busca
	s.HandleFunc(baseURL+users+filter, newHandlerFilter.Filter).Methods(http.MethodGet)

	// Rotas por ID
	s.HandleFunc(baseURL+user+idPath, newHandlerUser.GetByID).Methods(http.MethodGet)
	s.HandleFunc(baseURL+user+version+idParam, newHandlerUser.GetVersionByID).Methods(http.MethodGet)
	s.HandleFunc(baseURL+user+idParam, newHandlerUser.Update).Methods(http.MethodPut)
	s.HandleFunc(baseURL+user+idParam, newHandlerUser.Delete).Methods(http.MethodDelete)

	// Rotas de status
	s.HandleFunc(baseURL+user+enable+idParam, newHandlerUser.Enable).Methods(http.MethodPatch)
	s.HandleFunc(baseURL+user+disable+idParam, newHandlerUser.Disable).Methods(http.MethodPatch)
}
