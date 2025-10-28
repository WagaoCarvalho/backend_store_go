package routes

import (
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/config"
	handler "github.com/WagaoCarvalho/backend_store_go/internal/handler/product/category"
	jwtAuth "github.com/WagaoCarvalho/backend_store_go/internal/pkg/auth/jwt"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	jwt "github.com/WagaoCarvalho/backend_store_go/internal/pkg/middleware/jwt"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterProductCategoryRoutes(
	r *mux.Router,
	db *pgxpool.Pool,
	log *logger.LogAdapter,
	blacklist jwt.TokenBlacklist,
) {
	productCategoryRepo := repo.NewProductCategory(db)
	productCategoryService := service.NewProductCategory(productCategoryRepo)
	productCategoryHandler := handler.NewProductCategory(productCategoryService, log)

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

	s.HandleFunc("/product-category", productCategoryHandler.Create).Methods(http.MethodPost)
	s.HandleFunc("/product-category/{id:[0-9]+}", productCategoryHandler.GetByID).Methods(http.MethodGet)
	s.HandleFunc("/product-categories", productCategoryHandler.GetAll).Methods(http.MethodGet)
	s.HandleFunc("/product-category/{id:[0-9]+}", productCategoryHandler.Update).Methods(http.MethodPut)
	s.HandleFunc("/product-category/{id:[0-9]+}", productCategoryHandler.Delete).Methods(http.MethodDelete)
}
