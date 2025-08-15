package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/supplier_categories"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier_categories"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_categories"
	jwt_auth "github.com/WagaoCarvalho/backend_store_go/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/pkg/middleware/jwt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierCategoryRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LoggerAdapter,
	blacklist jwt.TokenBlacklist,
) {
	supplierCategoryRepo := repo.NewSupplierCategoryRepository(db, log)
	supplierCategoryService := service.NewSupplierCategoryService(supplierCategoryRepo, log)
	supplierCategoryHandler := handler.NewSupplierCategoryHandler(supplierCategoryService, log)

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
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager)) // <- passe o jwtManager, não a string SecretKey

	s.HandleFunc("/supplier-category", supplierCategoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/supplier-categories", supplierCategoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.Delete).Methods(http.MethodDelete)
}
