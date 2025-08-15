package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/supplier_category_relations"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_category_relations"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_category_relations"
	jwt_auth "github.com/WagaoCarvalho/backend_store_go/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/pkg/middleware/jwt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierCategoryRelationRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	relationRepo := repo.NewSupplierCategoryRelationRepo(db, log)
	relationService := service.NewSupplierCategoryRelationService(relationRepo, log)
	relationHandler := handler.NewSupplierCategoryRelationHandler(relationService, log)

	// Carregar config JWT
	jwtCfg := config.LoadJwtConfig()

	// Criar jwtManager que implementa JWTService
	jwtManager := jwt_auth.NewJWTManager(
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
