package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/category"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductCategoryRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL
	idPath := serverConfig.IDPath

	// Repositórios
	newRepoCategory := repo.NewProductCategory(db)

	// Serviços
	newServiceCategory := service.NewProductCategoryService(newRepoCategory)

	// Handlers
	newHandlerCategory := handler.NewProductCategoryHandler(newServiceCategory, log)

	// Config JWT
	jwtCfg := config.LoadJwtConfig()
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	// Rotas protegidas
	s := r.PathPrefix("/").Subrouter()
	s.Use(jwtMiddlewares.IsAuthByBearerToken(blacklist, log, jwtManager))

	// Constantes para caminhos
	const (
		productCategory   = "/product-category"
		productCategories = "/product-categories"
	)

	// Rotas de categoria de produto
	s.HandleFunc(baseURL+productCategory, newHandlerCategory.Create).Methods(http.MethodPost)
	s.HandleFunc(baseURL+productCategories, newHandlerCategory.GetAll).Methods(http.MethodGet)
	s.HandleFunc(baseURL+productCategory+idPath, newHandlerCategory.GetByID).Methods(http.MethodGet)
	s.HandleFunc(baseURL+productCategory+idPath, newHandlerCategory.Update).Methods(http.MethodPut)
	s.HandleFunc(baseURL+productCategory+idPath, newHandlerCategory.Delete).Methods(http.MethodDelete)
}
