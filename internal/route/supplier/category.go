package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/supplier/category"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterSupplierCategoryRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	supplierCategoryRepo := repo.NewSupplierCategory(db)
	supplierCategoryService := service.NewSupplierCategory(supplierCategoryRepo)
	supplierCategoryHandler := handler.NewSupplierCategoryHandler(supplierCategoryService, log)

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
	s.Use(jwt.IsAuthByBearerToken(blacklist, log, jwtManager)) // <- passe o jwtManager, não a string SecretKey

	s.HandleFunc("/supplier-category", supplierCategoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/supplier-categories", supplierCategoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/supplier-category/{id:[0-9]+}", supplierCategoryHandler.Delete).Methods(http.MethodDelete)
}
