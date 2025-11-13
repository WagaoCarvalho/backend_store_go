package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/category_relation"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category_relation"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category_relation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	relationRepo := repo.NewProductCategoryRelation(db)
	relationService := service.NewProductCategoryRelation(relationRepo)
	relationHandler := handler.NewProductCategoryRelationHandler(relationService, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Criar jwtManager que implementa JWTService
	jwtManager := jwtAuth.NewJWTManager(
		jwtCfg.SecretKey,
		jwtCfg.TokenDuration,
		jwtCfg.Issuer,
		jwtCfg.Audience,
	)

	s := r.PathPrefix("/").Subrouter()
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager))

	s.HandleFunc("/product-category-relation", relationHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/product-category-relations/{product_id:[0-9]+}", relationHandler.GetAllRelationsByProductID).Methods(http.MethodGet)
	s.HandleFunc("/product-category-relation/{product_id:[0-9]+}/category/{category_id:[0-9]+}", relationHandler.Delete).Methods(http.MethodDelete)
	s.HandleFunc("/product-category-relation/{product_id:[0-9]+}", relationHandler.DeleteAll).Methods(http.MethodDelete)
}
