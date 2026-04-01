package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/category_relation"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwtMiddlewares "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category_relation"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category_relation"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwtMiddlewares.TokenBlacklist,
) {
	serverConfig := config.LoadServerConfig()
	baseURL := serverConfig.BaseURL

	newRepoRelation := repo.NewProductCategoryRelation(db)
	newServiceRelation := service.NewProductCategoryRelation(newRepoRelation)
	newHandlerRelation := handler.NewProductCategoryRelationHandler(newServiceRelation, log)

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
		product                  = "/product"
		productCategoryRelation  = "/product-category-relation"
		productCategoryRelations = "/product-category-relations"
		productIDParam           = "/{product_id:[0-9]+}"
		categoryIDParam          = "/{category_id:[0-9]+}"
		category                 = "/category"
	)

	s.HandleFunc(baseURL+product+productIDParam+productCategoryRelations, newHandlerRelation.GetAllRelationsByProductID).Methods(http.MethodGet)
	s.HandleFunc(baseURL+product+productIDParam+productCategoryRelations, newHandlerRelation.Create).Methods(http.MethodPost)
	s.HandleFunc(baseURL+product+productIDParam+category+categoryIDParam, newHandlerRelation.Delete).Methods(http.MethodDelete)
	s.HandleFunc(baseURL+product+productIDParam+productCategoryRelations, newHandlerRelation.DeleteAll).Methods(http.MethodDelete)
}
