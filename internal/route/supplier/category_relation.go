package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/category_relation"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category_relation"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category_relation"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	relationRepo := repo.NewSupplierCategoryRelation(db)
	relationService := service.NewSupplierCategoryRelation(relationRepo)
	relationHandler := handler.NewSupplierCategoryRelation(relationService, log)

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

	s.HandleFunc("/supplier-category-relation", relationHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/supplier-category-relations/{supplier_id:[0-9]+}", relationHandler.GetBySupplierID).Methods(http.MethodGet)
	s.HandleFunc("/supplier-category-relation/{supplier_id:[0-9]+}/category/{category_id:[0-9]+}/exists", relationHandler.HasSupplierCategoryRelation).Methods(http.MethodGet)
	s.HandleFunc("/supplier-category-relation/{supplier_id:[0-9]+}/category/{category_id:[0-9]+}", relationHandler.DeleteByID).Methods(http.MethodDelete)
	s.HandleFunc("/supplier-category-relation/{supplier_id:[0-9]+}", relationHandler.DeleteAllBySupplierID).Methods(http.MethodDelete)
}
